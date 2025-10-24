ALTER TABLE urls
ALTER COLUMN last_checked TYPE timestamp USING last_checked AT TIME ZONE 'UTC',
ALTER COLUMN created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC';