-- Verify junction:appschema on pg

BEGIN;

SELECT pg_catalog.has_schema_privilege('junction', 'usage');

ROLLBACK;
