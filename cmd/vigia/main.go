package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	accountapp "github.com/vigia/vigia-v1/internal/account/application"
	accountpostgres "github.com/vigia/vigia-v1/internal/account/infrastructure/persistence/postgres"
	accounthttp "github.com/vigia/vigia-v1/internal/account/interfaces/http"
	notifapp "github.com/vigia/vigia-v1/internal/notification/application"
	notifobservability "github.com/vigia/vigia-v1/internal/notification/infrastructure/observability"
	notifpostgres "github.com/vigia/vigia-v1/internal/notification/infrastructure/persistence/postgres"
	"github.com/vigia/vigia-v1/internal/notification/infrastructure/provider/uazapi"
	"github.com/vigia/vigia-v1/internal/observability/analyzer"
	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/persistence/postgres"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/scheduler"
	httpapi "github.com/vigia/vigia-v1/internal/observability/interfaces/http"
	"github.com/vigia/vigia-v1/internal/shared/clock"
	"github.com/vigia/vigia-v1/internal/shared/id"
	"github.com/vigia/vigia-v1/internal/shared/middleware"
)

func main() {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, envOrDefault("DATABASE_URL", "postgres://vigia:vigia@localhost:5432/vigia"))
	if err != nil {
		log.Fatalf("vigia: failed to connect to database: %v", err)
	}
	defer pool.Close()

	monitors := postgres.NewMonitorRepository(pool)
	incidents := postgres.NewIncidentRepository(pool)
	samples := postgres.NewSampleRepository(pool)

	sysClock := clock.System()
	ids := id.Random()

	// Account context
	accountRepo := accountpostgres.NewAccountRepository(pool)
	getAccount := accountapp.NewGetAccount(accountRepo, sysClock)
	updateAccount := accountapp.NewUpdateAccount(accountRepo, sysClock)
	resolveRecipient := accountapp.NewResolveRecipient(accountRepo)

	// Notification context
	notifRepo := notifpostgres.NewNotificationRepository(pool)
	whatsappProvider := uazapi.NewWhatsAppProvider(
		envOrDefault("UAZAPI_BASE_URL", ""),
		envOrDefault("UAZAPI_TOKEN", ""),
		nil,
	)
	enqueueNotification := notifapp.NewEnqueueNotification(notifRepo, ids, sysClock)
	deliverNotifications := notifapp.NewDeliverNotifications(notifRepo, whatsappProvider, sysClock)
	publisher := notifobservability.NewPublisher(enqueueNotification, resolveRecipient)

	// Observability context
	resolveIncident := application.NewResolveIncident(incidents, publisher, sysClock)
	checkMonitor := application.NewCheckMonitor(
		monitors, incidents, samples,
		collector.NewHTTPCollector(nil), analyzer.NewThresholdAnalyzer(),
		resolveIncident, publisher, sysClock, ids,
	)
	buildDailySummary := application.NewBuildDailySummary(monitors, incidents, samples, publisher, sysClock)

	createMonitor := application.NewCreateMonitor(monitors, ids)
	pauseMonitor := application.NewPauseMonitor(monitors)
	resumeMonitor := application.NewResumeMonitor(monitors)
	disableMonitor := application.NewDisableMonitor(monitors)
	queryHistory := application.NewQueryHistory(incidents, samples)
	queryMonitors := application.NewQueryMonitors(monitors, incidents, samples)
	queryIncidents := application.NewQueryIncidents(incidents, monitors)
	queryAggregateHistory := application.NewQueryAggregateHistory(monitors, incidents, sysClock)

	sched := scheduler.NewTickerScheduler(monitors, checkMonitor, 30*time.Second, sysClock, nil)
	go sched.Run(ctx)
	go runDailySummaryLoop(ctx, buildDailySummary)
	go runDeliveryLoop(ctx, deliverNotifications)

	// HTTP
	supabaseURL := envOrDefault("SUPABASE_URL", "")
	jwksURL := supabaseURL + "/auth/v1/.well-known/jwks.json"
	authMW, err := middleware.NewAuth(jwksURL)
	if err != nil {
		log.Fatalf("vigia: failed to initialize auth middleware: %v", err)
	}

	obsHandlers := httpapi.NewHandlers(createMonitor, pauseMonitor, resumeMonitor, disableMonitor, queryHistory, queryMonitors, queryIncidents, queryAggregateHistory)
	accHandlers := accounthttp.NewHandlers(getAccount, updateAccount)

	mux := http.NewServeMux()
	mux.Handle("/account", authMW(accounthttp.NewRouter(accHandlers)))
	mux.Handle("/account/", authMW(accounthttp.NewRouter(accHandlers)))
	mux.Handle("/", authMW(httpapi.NewRouter(obsHandlers)))

	addr := envOrDefault("HTTP_ADDR", ":8080")
	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		log.Printf("vigia: HTTP listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("vigia: HTTP server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("vigia: shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("vigia: graceful shutdown failed: %v", err)
	}
}

func runDailySummaryLoop(ctx context.Context, uc *application.BuildDailySummary) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := uc.Execute(ctx); err != nil {
				log.Printf("vigia: daily summary failed: %v", err)
			}
		}
	}
}

func runDeliveryLoop(ctx context.Context, uc *notifapp.DeliverNotifications) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := uc.Execute(ctx); err != nil {
				log.Printf("vigia: notification delivery failed: %v", err)
			}
		}
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
