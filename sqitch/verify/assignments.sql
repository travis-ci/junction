-- Verify junction:assignments on pg

BEGIN;

SELECT id, worker_id, job_id, dispatched, last_heartbeat FROM junction.assignments WHERE false;

ROLLBACK;
