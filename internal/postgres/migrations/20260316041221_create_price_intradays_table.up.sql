BEGIN;

CREATE TABLE IF NOT EXISTS price_intradays (
    symbol VARCHAR(4) NOT NULL,
    unix_timestamp BIGINT NOT NULL,
    open INT NOT NULL,
    close INT NOT NULL,
    high INT NOT NULL,
    low INT NOT NULL,
    transaction_value BIGINT NOT NULL,
    volume BIGINT NOT NULL,
    PRIMARY KEY (symbol, unix_timestamp)
);

COMMIT;
