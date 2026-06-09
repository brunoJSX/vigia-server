-- Schema for the Observability context's Postgres persistence.
-- Plain SQL, no migration tool — explicit and small enough to read in full.

CREATE TABLE IF NOT EXISTS monitors (
    id          TEXT PRIMARY KEY,
    account_id  TEXT NOT NULL,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    target      TEXT NOT NULL,
    type        TEXT NOT NULL,
    status      TEXT NOT NULL,
    threshold                    INTEGER NOT NULL,
    interval_ns                  BIGINT NOT NULL,
    acceptable_response_time_ns  BIGINT,
    created_at                   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS incidents (
    id              TEXT PRIMARY KEY,
    monitor_id      TEXT NOT NULL REFERENCES monitors (id),
    status          TEXT NOT NULL,
    opened_at       TIMESTAMPTZ NOT NULL,
    resolved_at     TIMESTAMPTZ,
    sequence_number SERIAL
);

CREATE INDEX IF NOT EXISTS incidents_monitor_status_idx ON incidents (monitor_id, status);
CREATE INDEX IF NOT EXISTS incidents_monitor_opened_at_idx ON incidents (monitor_id, opened_at);

CREATE TABLE IF NOT EXISTS samples (
    monitor_id TEXT NOT NULL REFERENCES monitors (id),
    "timestamp" TIMESTAMPTZ NOT NULL,
    success    BOOLEAN NOT NULL,
    latency_ns BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS samples_monitor_timestamp_idx ON samples (monitor_id, "timestamp");
