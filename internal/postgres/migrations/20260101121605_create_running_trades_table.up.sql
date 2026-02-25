BEGIN;

CREATE TABLE IF NOT EXISTS running_trades (
    symbol VARCHAR(4) NOT NULL,
    date DATE NOT NULL,
    buyer VARCHAR(2) NOT NULL,
    buyer_investor_type VARCHAR(1) NOT NULL,
    seller VARCHAR(2) NOT NULL,
    seller_investor_type VARCHAR(1) NOT NULL,
    market_board VARCHAR(2) NOT NULL,
    action VARCHAR(1) NOT NULL,
    price INT NOT NULL,
    lot INT NOT NULL,
    trade_number BIGINT NOT NULL,
    PRIMARY KEY (symbol, date, trade_number)
);

COMMENT ON COLUMN running_trades.action IS 'B: Buy, S: Sell';
COMMENT ON COLUMN running_trades.buyer_investor_type IS 'D: Domestic, F: Foreign';
COMMENT ON COLUMN running_trades.seller_investor_type IS 'D: Domestic, F: Foreign';
COMMENT ON COLUMN running_trades.market_board IS 'RG: Regular, NG: Negosiasi, TN: Tunai';

COMMIT;
