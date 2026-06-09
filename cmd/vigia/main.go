// Bootstrap for the Observability vertical slice — manual dependency
// injection, no framework (golang-conventions: prefer explicit code).
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

	"github.com/vigia/vigia-v1/internal/observability/analyzer"
	"github.com/vigia/vigia-v1/internal/observability/application"
	"github.com/vigia/vigia-v1/internal/observability/collector"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/notification"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/persistence/postgres"
	"github.com/vigia/vigia-v1/internal/observability/infrastructure/scheduler"
	httpapi "github.com/vigia/vigia-v1/internal/observability/interfaces/http"
	"github.com/vigia/vigia-v1/internal/shared/clock"
	"github.com/vigia/vigia-v1/internal/shared/id"
)

func main() {
	// Loads .env into the process environment when present — silently a
	// no-op otherwise, so real deployments (env vars set directly) work
	// the same way.
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

	publisher := notification.NewStubPublisher(nil)
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

	handlers := httpapi.NewHandlers(createMonitor, pauseMonitor, resumeMonitor, disableMonitor, queryHistory, queryMonitors, queryIncidents, queryAggregateHistory)
	addr := envOrDefault("HTTP_ADDR", ":8080")
	server := &http.Server{Addr: addr, Handler: httpapi.NewRouter(handlers)}

	go func() {
		log.Printf("vigia: observability HTTP listening on %s", addr)
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

// runDailySummaryLoop fires BuildDailySummary roughly once a day — the
// "execução diária agendada" the workflow describes. Exact wall-clock timing
// (e.g. midnight) is an infra detail left for whenever it actually matters.
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

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
