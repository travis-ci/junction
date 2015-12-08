-- Revert junction:assignments from pg

BEGIN;

DROP TABLE junction.assignments;

COMMIT;
