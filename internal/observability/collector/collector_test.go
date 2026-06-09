package collector_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vigia/vigia-v1/internal/observability/collector"
)

func serverWithStatus(code int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	}))
}

func TestCollect_2xx_IsSuccess(t *testing.T) {
	srv := serverWithStatus(http.StatusOK)
	defer srv.Close()

	c := collector.NewHTTPCollector(nil)
	s, err := c.Collect(context.Background(), srv.URL, 5*time.Second)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Success {
		t.Fatal("expected 200 to be Success=true")
	}
}

func TestCollect_3xx_IsSuccess(t *testing.T) {
	srv := serverWithStatus(http.StatusFound)
	defer srv.Close()

	c := collector.NewHTTPCollector(nil)
	s, err := c.Collect(context.Background(), srv.URL, 5*time.Second)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Success {
		t.Fatal("expected 302 to be Success=true")
	}
}

// CD-001: 4xx counts as failure.
func TestCollect_4xx_IsFailure(t *testing.T) {
	srv := serverWithStatus(http.StatusNotFound)
	defer srv.Close()

	c := collector.NewHTTPCollector(nil)
	s, err := c.Collect(context.Background(), srv.URL, 5*time.Second)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Success {
		t.Fatal("expected 404 to be Success=false (CD-001)")
	}
}

// CD-001: 5xx counts as failure.
func TestCollect_5xx_IsFailure(t *testing.T) {
	srv := serverWithStatus(http.StatusInternalServerError)
	defer srv.Close()

	c := collector.NewHTTPCollector(nil)
	s, err := c.Collect(context.Background(), srv.URL, 5*time.Second)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Success {
		t.Fatal("expected 500 to be Success=false (CD-001)")
	}
}

func TestCollect_Timeout_IsFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := collector.NewHTTPCollector(nil)
	s, err := c.Collect(context.Background(), srv.URL, 50*time.Millisecond)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Success {
		t.Fatal("expected timeout to produce Success=false")
	}
}
