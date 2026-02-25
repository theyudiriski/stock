BEGIN;

CREATE TABLE IF NOT EXISTS emitten_subsidiaries (
    symbol VARCHAR(4) NOT NULL,
    subsidiary_name VARCHAR(255) NOT NULL,
    subsidiary_percentage DECIMAL(10, 6) NOT NULL,
    subsidiary_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (symbol, subsidiary_name)
);

COMMIT;
