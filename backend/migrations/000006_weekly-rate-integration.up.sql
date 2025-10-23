CREATE TABLE weekly_rate (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    rate_name VARCHAR(255) NOT NULL,
    amount SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users
ADD COLUMN weekly_rate_id INT,
ADD CONSTRAINT fk_weekly_rate FOREIGN KEY (weekly_rate_id) REFERENCES weekly_rate (id);

CREATE INDEX idx_users_weekly_rate_id ON users (weekly_rate_id);

INSERT INTO weekly_rate (uuid, rate_name, amount)
VALUES (gen_random_uuid()::varchar, 'Temps pleins', 35);

INSERT INTO weekly_rate (uuid, rate_name, amount)
VALUES (gen_random_uuid()::varchar, 'Temps pleins + RTT', 39);

INSERT INTO weekly_rate (uuid, rate_name, amount)
VALUES (gen_random_uuid()::varchar, 'Temps partiel', 20);