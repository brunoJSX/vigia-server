package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	accountapp "github.com/vigia/vigia-v1/internal/account/application"
	"github.com/vigia/vigia-v1/internal/shared/middleware"
)

type Handlers struct {
	getAccount    *accountapp.GetAccount
	updateAccount *accountapp.UpdateAccount
}

func NewHandlers(getAccount *accountapp.GetAccount, updateAccount *accountapp.UpdateAccount) *Handlers {
	return &Handlers{getAccount: getAccount, updateAccount: updateAccount}
}

type accountResponse struct {
	ID             string    `json:"id"`
	WhatsAppNumber *string   `json:"whatsapp_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (h *Handlers) GetAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	a, err := h.getAccount.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toResponse(a.ID, a.WhatsAppNumber, a.CreatedAt, a.UpdatedAt))
}

type updateAccountRequest struct {
	WhatsAppNumber *string `json:"whatsapp_number"`
}

func (h *Handlers) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req updateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	number := ""
	if req.WhatsAppNumber != nil {
		number = *req.WhatsAppNumber
	}

	a, err := h.updateAccount.Execute(r.Context(), userID, accountapp.UpdateAccountInput{
		WhatsAppNumber: number,
	})
	if err != nil {
		if errors.Is(err, accountapp.ErrInvalidPhoneNumber) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, accountapp.ErrAccountNotFound) {
			http.Error(w, "account not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toResponse(a.ID, a.WhatsAppNumber, a.CreatedAt, a.UpdatedAt))
}

func toResponse(id, whatsappNumber string, createdAt, updatedAt time.Time) accountResponse {
	resp := accountResponse{ID: id, CreatedAt: createdAt, UpdatedAt: updatedAt}
	if whatsappNumber != "" {
		resp.WhatsAppNumber = &whatsappNumber
	}
	return resp
}
