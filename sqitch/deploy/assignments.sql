-- Deploy junction:assignments to pg
-- requires: workers
-- requires: appschema

BEGIN;

CREATE TABLE junction.assignments (
    id         uuid    PRIMARY KEY,
    worker_id  uuid    NOT NULL REFERENCES junction.workers,
    job_id     bigint  NOT NULL,
	dispatched boolean NOT NULL,
	last_heartbeat timestamptz
);

COMMIT;
