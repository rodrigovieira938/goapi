CREATE TABLE IF NOT EXISTS car (
    id             SERIAL PRIMARY KEY,
    model          TEXT NOT NULL,
    brand          TEXT NOT NULL,
    year           INT NOT NULL,
    color          TEXT,
    doors          INT,
    price_per_day  NUMERIC(10,2),   -- for currency
    license_plate  TEXT UNIQUE,
    baggage_volume NUMERIC(10,2)    -- liters
);
