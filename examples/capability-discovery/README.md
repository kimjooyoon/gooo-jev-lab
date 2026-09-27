# Capability discovery fixture

This declaration is a small, source-bound example for asking what `gooo` can currently do. It is intentionally not an execution request and it does not grant authorization.

Run it from a checkout of the canonical [`gooo-jev`](https://github.com/kimjooyoon/gooo-jev) repository that contains the capability-discovery implementation:

```sh
go run ./cmd/gooo-discover \
  --query "What can gooo do with this declaration?" \
  /path/to/gooo-jev-lab/examples/capability-discovery/declaration.gooo
```

The response should expose a capability status, source digest, and observed declaration signals. A digest or signal is provenance evidence only; it is not proof that a capability executed, changed a repository, or crossed an authorization boundary.
