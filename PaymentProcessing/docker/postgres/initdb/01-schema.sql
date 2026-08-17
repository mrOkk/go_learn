CREATE TABLE IF NOT EXISTS merchants (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    balance NUMERIC(18, 2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    amount NUMERIC(18, 2) NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('approved', 'rejected')),
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
