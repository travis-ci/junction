-- Verify junction:workers on pg

BEGIN;

SELECT id, queue, created_at, last_heartbeat, max_job_count, attributes FROM junction.workers WHERE false;

ROLLBACK;
