CREATE TYPE valid_status AS ENUM (
    'new',
    'success',
    'failure',
    'error',
    'canceled'
);

CREATE TYPE valid_currency AS ENUM (
    'usd',
    'eur',
    'rub'
);

CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    user_email VARCHAR(20) NOT NULL,
    currency valid_currency NOT NULL,
    amount decimal(12, 2) NOT NULL CHECK(amount > 0) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status valid_status NOT NULL DEFAULT 'new'
);

CREATE INDEX ON payments(user_email);
CREATE INDEX ON payments(user_id);

CREATE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
    BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();