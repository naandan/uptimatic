CREATE TABLE status_logs (
    id SERIAL PRIMARY KEY,
    url_id INT REFERENCES urls(id),
    status VARCHAR(20) NOT NULL,
    response_time INT,
    checked_at TIMESTAMP DEFAULT NOW()
);