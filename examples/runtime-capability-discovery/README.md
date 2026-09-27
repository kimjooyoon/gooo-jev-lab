# Runtime capability discovery fixture

This fixture demonstrates the runtime-side, source-bound capability query added to `gooo-jev-runtime`. It observes the declaration bytes and preserves a digest and structural signals; it does not execute work or authorize a capability.

Run from a checkout of the canonical runtime repository:

```sh
go run ./cmd/jevcal-discover < /path/to/gooo-jev-lab/examples/runtime-capability-discovery/request.json
```

The output is an observation surface. `declaration_source_digest` binds the input bytes, while `declaration_observed_signals` reports only structural markers. Neither field is semantic completion evidence or permission.
