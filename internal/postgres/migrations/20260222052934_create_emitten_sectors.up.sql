BEGIN;

CREATE TABLE IF NOT EXISTS emitten_sectors (
    id INT NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INT DEFAULT NULL
);

INSERT INTO emitten_sectors (id, name, parent_id) VALUES
(1, 'Barang Konsumen Primer', NULL),
(2, 'Kesehatan', NULL),
(3, 'Keuangan', NULL),
(4, 'Barang Konsumen Non-Primer', NULL),
(5, 'Properti & Real Estat', NULL),
(6, 'Perindustrian', NULL),
(7, 'Energi', NULL),
(8, 'Barang Baku', NULL),
(9, 'Infrastruktur', NULL),
(50, 'Teknologi', NULL),
(51, 'Transportasi & Logistik', NULL),
(65, 'Currencies', NULL),
(67, 'Others', NULL),
(70, 'Indeks Sektoral', NULL),
(73, 'Commodities', NULL),
(75, 'Cryptocurrency', NULL),
(78, 'Global Index', NULL),
(80, 'Reksadana', NULL),
(85, 'Delisted Stock', NULL),
(88, 'Indeks', NULL),
(89, 'Listing Board', NULL),
(92, 'FR Bonds', NULL);

COMMIT;
