BEGIN;

CREATE TABLE IF NOT EXISTS broker_summaries (
    symbol VARCHAR(4) NOT NULL,
    broker VARCHAR(2) NOT NULL,
    action VARCHAR(1) NOT NULL,
    investor_type VARCHAR(1) NOT NULL,
    market_board VARCHAR(2) NOT NULL,
    summary_date DATE NOT NULL,
    total_lot BIGINT NOT NULL,
    total_value BIGINT NOT NULL,
    price_average NUMERIC(20, 2) NOT NULL,
    PRIMARY KEY (symbol, broker, action, investor_type, market_board, summary_date)
);

COMMENT ON COLUMN broker_summaries.action IS 'B: Buy, S: Sell';
COMMENT ON COLUMN broker_summaries.investor_type IS 'D: Domestic, F: Foreign';
COMMENT ON COLUMN broker_summaries.market_board IS 'RG: Regular, NG: Negosiasi, TN: Tunai';

COMMIT;
