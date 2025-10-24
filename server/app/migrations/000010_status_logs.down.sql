ALTER TABLE status_logs
ALTER COLUMN checked_at TYPE timestamp USING checked_at AT TIME ZONE 'UTC';