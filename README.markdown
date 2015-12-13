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
