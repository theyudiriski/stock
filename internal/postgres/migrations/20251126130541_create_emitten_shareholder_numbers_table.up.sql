BEGIN;

CREATE TABLE IF NOT EXISTS emitten_shareholder_numbers (
    symbol VARCHAR(4) NOT NULL,
    shareholder_date DATE NOT NULL,
    total_share BIGINT NOT NULL,
    change INT NOT NULL,
    PRIMARY KEY (symbol, shareholder_date)
);

COMMIT;
