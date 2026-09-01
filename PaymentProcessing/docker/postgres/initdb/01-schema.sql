CREATE TABLE IF NOT EXISTS merchants (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    balance NUMERIC(18, 2) NOT NULL DEFAULT 0
);

INSERT INTO merchants (id, name, is_active, balance)
SELECT v.id, v.name, v.is_active, round((random() * 100000)::numeric, 2)
FROM (VALUES
    (1, 'Sunrise Bakery', TRUE),
    (2, 'Blue Horizon Electronics', TRUE),
    (3, 'Green Leaf Grocers', TRUE),
    (4, 'Old Town Bookstore', FALSE),
    (5, 'Silver Line Motors', FALSE)
) AS v(id, name, is_active)
WHERE NOT EXISTS (SELECT 1 FROM merchants);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES merchants(id),
    amount NUMERIC(18, 2) NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('approved', 'rejected')),
    "timestamp" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
