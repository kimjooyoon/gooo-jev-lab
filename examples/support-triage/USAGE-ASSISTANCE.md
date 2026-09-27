# Usage assistance fixture

This fixture describes how a `.gooo` declaration can discover what the language can currently support without executing a task or granting authorization.

The canonical implementation is in [gooo-jev](https://github.com/kimjooyoon/gooo-jev):

```sh
go run ./cmd/gooo-assist examples/support-triage/partial.gooo pro
go run ./cmd/gooo-assist --json examples/support-triage/partial.gooo pro
```

The contract keeps three states distinct:

- `AVAILABLE`: the capability is bound to the current declaration evidence.
- `DEFERRED`: the next operation is known, but the required evidence boundary is not bound.
- `UNKNOWN`: the declaration is not sufficient to make a capability claim.

Every structured result carries source, discovery, plan, and assistance digests. The fixture is non-executing and non-authorizing; a DEFERRED result is neither success nor failure.
