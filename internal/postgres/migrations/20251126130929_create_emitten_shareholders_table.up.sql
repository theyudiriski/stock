BEGIN;

CREATE TABLE IF NOT EXISTS emitten_shareholders (
    symbol VARCHAR(4) NOT NULL,
    shareholder_name TEXT NOT NULL,
    shareholder_percentage NUMERIC(10, 6) NOT NULL,
    shareholder_badges TEXT[] DEFAULT '{}',
    PRIMARY KEY (symbol, shareholder_name)
);

COMMIT;
