BEGIN;

CREATE TABLE IF NOT EXISTS emitten_subsidiary_companies (
    symbol VARCHAR(4) NOT NULL,
    subsidiary_company_name VARCHAR(255) NOT NULL,
    subsidiary_company_percentage DECIMAL(10, 6) DEFAULT NULL,
    subsidiary_company_type TEXT NOT NULL,
    PRIMARY KEY (symbol, subsidiary_company_name)
);


CREATE TABLE IF NOT EXISTS emitten_subsidiary_companies_history (
    symbol VARCHAR(4) NOT NULL,
    last_updated_period VARCHAR(7) NOT NULL,
    PRIMARY KEY (symbol)
);

COMMIT;
