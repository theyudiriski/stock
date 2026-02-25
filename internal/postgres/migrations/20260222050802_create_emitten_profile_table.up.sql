BEGIN;

CREATE TABLE IF NOT EXISTS emitten_profiles (
    symbol VARCHAR(4) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    underwriters TEXT[] NOT NULL,
    underwriters_code TEXT[] DEFAULT '{}',
    free_float DECIMAL(10, 6) DEFAULT NULL,
    subsector_id INT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (symbol)
);

COMMIT;
