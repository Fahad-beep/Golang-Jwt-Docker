CREATE TABLE IF NOT EXISTS otp (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    code INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);