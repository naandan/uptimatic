CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    label VARCHAR(255),
    url TEXT NOT NULL,
    last_checked TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);