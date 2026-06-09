CREATE TABLE IF NOT EXISTS notifications (
    id           TEXT PRIMARY KEY,
    type         TEXT NOT NULL,
    recipient    TEXT NOT NULL,
    payload      JSONB NOT NULL DEFAULT '{}',
    status       TEXT NOT NULL DEFAULT 'pending',
    attempts     INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL,
    delivered_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS notifications_actionable_idx
    ON notifications (status, attempts)
    WHERE status IN ('pending', 'failed');
