BEGIN;

CREATE TABLE IF NOT EXISTS investor_categories (
    category_code VARCHAR(2) NOT NULL,
    category_name VARCHAR(30) NOT NULL,
    PRIMARY KEY (category_code)
);

-- Source: https://www.idxchannel.com/market-news/simak-9-tipe-investor-dalam-klasifikasi-ksei-apa-saja-kategorinya/all
INSERT INTO investor_categories (category_code, category_name) VALUES ('ID', 'Individual');
INSERT INTO investor_categories (category_code, category_name) VALUES ('CP', 'Corporate');
INSERT INTO investor_categories (category_code, category_name) VALUES ('MF', 'Mutual Fund');
INSERT INTO investor_categories (category_code, category_name) VALUES ('IB', 'Financial Institution');
INSERT INTO investor_categories (category_code, category_name) VALUES ('IS', 'Insurance');
INSERT INTO investor_categories (category_code, category_name) VALUES ('SC', 'Securities Company');
INSERT INTO investor_categories (category_code, category_name) VALUES ('PF', 'Pension Fund');
INSERT INTO investor_categories (category_code, category_name) VALUES ('FD', 'Foundation');
INSERT INTO investor_categories (category_code, category_name) VALUES ('OT', 'Others');

CREATE TABLE IF NOT EXISTS emitten_shareholder_chart (
    symbol VARCHAR(4) NOT NULL,
    investor_category_codes VARCHAR(2)[] NOT NULL,
    investor_type VARCHAR(1) NOT NULL,
    date DATE NOT NULL,
    percentage DECIMAL(10, 6) NOT NULL,
    PRIMARY KEY (symbol, investor_category_codes, investor_type, date)
);

CREATE TABLE IF NOT EXISTS emitten_shareholder_chart_history (
    symbol VARCHAR(4) NOT NULL,
    last_updated_date VARCHAR(10) NOT NULL,
    PRIMARY KEY (symbol)
);

COMMIT;