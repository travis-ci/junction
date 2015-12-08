-- Revert junction:appschema from pg

BEGIN;

DROP SCHEMA junction;

COMMIT;
