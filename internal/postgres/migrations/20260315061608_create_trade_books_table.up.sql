BEGIN;

CREATE TABLE IF NOT EXISTS trade_books (
    symbol VARCHAR(4) NOT NULL,
    datetime TIMESTAMP NOT NULL,
    buy_lot INT NOT NULL,
    buy_frequency INT NOT NULL,
    sell_lot INT NOT NULL,
    sell_frequency INT NOT NULL,
    PRIMARY KEY (symbol, datetime)
);

COMMIT;
