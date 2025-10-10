ALTER TABLE status_logs
DROP CONSTRAINT IF EXISTS status_logs_url_id_fkey,
ADD CONSTRAINT status_logs_url_id_fkey
    FOREIGN KEY (url_id)
    REFERENCES urls(id)
    ON DELETE NO ACTION;
