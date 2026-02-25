BEGIN;

CREATE TABLE IF NOT EXISTS price_feeds (
    symbol VARCHAR(4) NOT NULL,
    date DATE NOT NULL,
    open INT NOT NULL,
    close INT NOT NULL,
    high INT NOT NULL,
    low INT NOT NULL,
    average INT NOT NULL,
    value BIGINT NOT NULL,
    volume BIGINT NOT NULL,
    frequency BIGINT NOT NULL,
    net_foreign BIGINT NOT NULL,
    PRIMARY KEY (symbol, date)
);

ALTER TABLE price_feeds ALTER COLUMN close DROP NOT NULL;

COMMIT;
