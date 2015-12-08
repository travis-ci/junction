-- Deploy junction:workers to pg
-- requires: appschema

BEGIN;

CREATE TABLE junction.workers (
    id             uuid        PRIMARY KEY,
    queue          text        NOT NULL,
    created_at     timestamptz NOT NULL,
    max_job_count  integer     NOT NULL,
    last_heartbeat timestamptz
);

COMMIT;
