BEGIN;

CREATE TABLE IF NOT EXISTS emitten_beneficiaries (
    symbol VARCHAR(4) NOT NULL,
    beneficiary_name VARCHAR(255) NOT NULL,
    PRIMARY KEY (symbol, beneficiary_name)
);

COMMIT;