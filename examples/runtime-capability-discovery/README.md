# Runtime capability discovery fixture

This fixture demonstrates the runtime-side, source-bound capability query added to gooo-jev-runtime. It observes the declaration bytes and preserves a digest and structural signals; it does not execute work or authorize a capability.

Run from a checkout of the canonical runtime repository:

~~~sh
go run ./cmd/jevcal-discover < /path/to/gooo-jev-lab/examples/runtime-capability-discovery/request.json
~~~

The output is an observation surface. declaration_source_digest binds the input bytes, while declaration_observed_signals reports only structural markers. The queries.json corpus shows three safe boundary cases:

- overview asks what the declaration can expose and returns the read-only catalog.
- provenance asks where the declaration came from and returns the source and generation evidence route.
- execution asks for an external boundary and remains deferred rather than being treated as permission.

The capability response carries the next operation and follow-up queries so a client can continue from evidence instead of guessing. None of these fields is semantic completion evidence or an authorization grant.
