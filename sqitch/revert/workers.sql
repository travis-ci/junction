-- Revert junction:workers from pg

BEGIN;

DROP TABLE junction.workers;

COMMIT;
