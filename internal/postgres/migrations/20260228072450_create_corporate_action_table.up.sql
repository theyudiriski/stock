BEGIN;

CREATE TABLE IF NOT EXISTS corpaction_dividend (
	symbol VARCHAR(4) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    cum_date DATE NOT NULL,
    ex_date DATE NOT NULL,
    recording_date DATE NOT NULL,
    payment_date DATE NOT NULL,
    PRIMARY KEY (symbol, cum_date, ex_date, recording_date, payment_date)
);

CREATE TABLE IF NOT EXISTS corpaction_rups (
    symbol VARCHAR(4) NOT NULL,
    datetime TIMESTAMPTZ DEFAULT NULL,
    venue TEXT NOT NULL,
    eligible_date DATE NOT NULL,
    PRIMARY KEY (symbol, eligible_date)
);

CREATE TABLE IF NOT EXISTS corpaction_public_expose (
    symbol VARCHAR(4) NOT NULL,
    datetime TIMESTAMPTZ NOT NULL,
    venue TEXT NOT NULL,
    PRIMARY KEY (symbol, datetime)
);

CREATE TABLE IF NOT EXISTS corpaction_right_issue (
    symbol VARCHAR(4) NOT NULL,
    price INT NOT NULL,
    old BIGINT NOT NULL,
    new BIGINT NOT NULL,
    cum_date DATE NOT NULL,
    ex_date DATE NOT NULL,
    recording_date DATE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    PRIMARY KEY (symbol, cum_date)
);

CREATE TABLE IF NOT EXISTS corpaction_stock_split (
    symbol VARCHAR(4) NOT NULL,
    old BIGINT NOT NULL,
    new BIGINT NOT NULL,
    cum_date DATE NOT NULL,
    ex_date DATE NOT NULL,
    recording_date DATE NOT NULL,
    PRIMARY KEY (symbol, cum_date)
);

CREATE TABLE IF NOT EXISTS corpaction_reverse_split (
    symbol VARCHAR(4) NOT NULL,
    old BIGINT NOT NULL,
    new BIGINT NOT NULL,
    cum_date DATE NOT NULL,
    ex_date DATE NOT NULL,
    recording_date DATE NOT NULL,
    PRIMARY KEY (symbol, cum_date)
);

COMMIT;