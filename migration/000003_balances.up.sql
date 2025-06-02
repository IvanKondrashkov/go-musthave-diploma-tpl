CREATE TABLE IF NOT EXISTS balances (
    user_id UUID REFERENCES users(id) PRIMARY KEY,
    current NUMERIC(10, 2) DEFAULT 0,
    withdrawn NUMERIC(10, 2) DEFAULT 0
);