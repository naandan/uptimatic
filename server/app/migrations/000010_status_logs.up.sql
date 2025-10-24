ALTER TABLE status_logs
ALTER COLUMN checked_at TYPE timestamptz USING checked_at AT TIME ZONE 'UTC';