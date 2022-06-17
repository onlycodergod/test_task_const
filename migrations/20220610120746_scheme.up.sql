// Creating a type called valid_status that can only be one of the values in the list.
CREATE TYPE valid_status AS ENUM (
    'new',
    'success',
    'failure',
    'error',
    'canceled'
);

// Creating a type called valid_currency that can only be one of the values in the list.
CREATE TYPE valid_currency AS ENUM (
    'usd',
    'eur',
    'rub'
);

// Creating a table called payments with the following columns:
// - id: a serial primary key
// - user_id: an integer that cannot be null
// - user_email: a string that cannot be null
// - currency: a valid_currency that cannot be null
// - amount: a decimal that cannot be null and must be greater than 0
// - created_at: a timestamp with time zone that cannot be null and defaults to now
// - updated_at: a timestamp with time zone that cannot be null and defaults to now
// - status: a valid_status that cannot be null and defaults to 'new'

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

// The above code is creating a trigger that will update the updated_at column with the current time
// whenever a row is updated.

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