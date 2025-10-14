ALTER TABLE status_logs
ALTER COLUMN status TYPE INTEGER
USING status::integer;