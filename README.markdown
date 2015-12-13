# Junction

**WARNING** This is _very_ much a work in progress. Here be dragons, ghosts, skulls, skeletons, and all sorts of terrifying creatures.

## TODO

Plans and other things:

- Store the version number with workers, to make it easy to see which version is rolled out everywhere.
- Assignments should probably be renamed. They are more like "a single job
  run". Eventually the worker will only interact with the job for a short
  amount of time, but the "assignment" will contain heartbeat timestamps for
  the entire duration of the job.
- The worker heartbeat endpoint takes a list of assignment IDs and also
  triggers a heartbeat for each of those assignments. The worker will send the
  list of assignments it's currently working on.
- The heartbeat handlers just update the `last_heartbeat` timestamp.
- A process separate from the web server will regularly get the list of jobs that scheduler has scheduled (`state=queued`, probably), sort them by `queued_at` and split them up by queue. Then, for each queue it will try to enqueue as many jobs as it can, starting at the oldest job by `queued_at`. Enqueueing a job is done by creating an assigment and somehow pushing it to the worker. I have no good idea for how to "push" the job yet.
- Junction knows about the maximum number of queueable jobs per queue (the "capacity"). Each worker can also have a max count, but doesn't need to.

## Overview

Currently, Junction's only job is to take jobs that are queued by
[scheduler](https://github.com/travis-ci/travis-scheduler) and make sure that
they run on a worker. It's mostly replacing the RabbitMQ "builds.something"
queues that we have today, and is also responsible for reporting state updates
back to hub.

The long-term goal for Junction is that it handles all data between the part of
Travis that runs jobs and the rest of Travis.

When Scheduler enqueues a job, it sets state=queued and the queued_at attribute
on the job in the database. Junction will periodically query the database and
list all of the jobs that are state=queued. It then determines the "queue" for
each job, and groups them by queue. For each queue, it'll start at the job that
was queued the longest ago (using queued_at) and will attempt to find a worker
to run the job on. Once it finds one, it will create an assignment for the job.

Each worker regularly sends heartbeats to Junction, including the assignments
it is currently working on. Junction will then update the last_heartbeat time
for the worker itself and all assignments in the database. Junction will return
the list of assignments that the worker is supposed to be working on.

When the worker gets back the list of assignments, it will terminate any
assignments it's working on that is not in the list, and will start up a
processor for every assignment it is not yet running.

If the list that worker sends to Junction is missing an assignment that at
least one heartbeat has been sent for, but no "final" state update (i.e. "job
finished" or "requeue job") has been received, Junction should assume that the
job needs to be requeued. This means that Worker needs to make sure to send
state updates before removing the assignment from it's internal list that it
sends to Junction in heartbeats.
