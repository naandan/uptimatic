ALTER TABLE urls
ALTER COLUMN last_checked TYPE timestamptz USING last_checked AT TIME ZONE 'UTC',
ALTER COLUMN created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC';