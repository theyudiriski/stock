BEGIN;

CREATE TABLE IF NOT EXISTS holidays (
    date DATE NOT NULL PRIMARY KEY
);

INSERT INTO holidays (date) VALUES ('2026-01-01');
INSERT INTO holidays (date) VALUES ('2026-01-16');
INSERT INTO holidays (date) VALUES ('2026-02-16');
INSERT INTO holidays (date) VALUES ('2026-02-17');
INSERT INTO holidays (date) VALUES ('2026-03-18');
INSERT INTO holidays (date) VALUES ('2026-03-19');
INSERT INTO holidays (date) VALUES ('2026-03-20');
INSERT INTO holidays (date) VALUES ('2026-03-23');
INSERT INTO holidays (date) VALUES ('2026-03-24');
INSERT INTO holidays (date) VALUES ('2026-04-03');
INSERT INTO holidays (date) VALUES ('2026-05-01');
INSERT INTO holidays (date) VALUES ('2026-05-14');
INSERT INTO holidays (date) VALUES ('2026-05-15');
INSERT INTO holidays (date) VALUES ('2026-05-27');
INSERT INTO holidays (date) VALUES ('2026-05-28');
INSERT INTO holidays (date) VALUES ('2026-06-01');
INSERT INTO holidays (date) VALUES ('2026-06-16');
INSERT INTO holidays (date) VALUES ('2026-08-17');
INSERT INTO holidays (date) VALUES ('2026-08-25');
INSERT INTO holidays (date) VALUES ('2026-12-24');
INSERT INTO holidays (date) VALUES ('2026-12-25');
INSERT INTO holidays (date) VALUES ('2026-12-31');

COMMIT;
