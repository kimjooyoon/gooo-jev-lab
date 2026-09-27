# gooo-jev-lab

An experimental Go and `.gooo` lab for provider-neutral JEV-style typed decisions.

This repository is intentionally a thin integration boundary, not a second copy of the language or runtime. The canonical implementation remains in [`gooo-jev`](https://github.com/kimjooyoon/gooo-jev) and [`gooo-jev-runtime`](https://github.com/kimjooyoon/gooo-jev-runtime).

## What this lab fixes in place

- `Choice`, `Score`, and `Noul` signals are represented as typed values rather than generated prose.
- Every receipt is bound to request/state identity and a declaration, IR, generation, and reverse-observation evidence chain.
- Invalid digests, distributions, confidence values, and question bindings fail closed.
- Confidence can select `ACCEPT`, `REVIEW`, or `REJECT` as an observation-only route. It never grants execution or authorization.
- The `.gooo` file is an experimental contract sketch until the canonical grammar promotes it.

## Why this exists

JEV-like models are useful when a program needs a bounded judgment, a probability distribution, or a score. The language should keep control flow, side effects, permissions, and provenance in code while using a model only as a typed observation source. This makes uncertainty visible and lets a human review boundary remain explicit.

## Validation

The repository validates in GitHub Actions only. Local Go execution is intentionally not part of the development contract.
## Capability discovery fixture

[`examples/capability-discovery`](examples/capability-discovery) provides a source-bound `.gooo` declaration and a canonical command for asking what the language can currently expose. It records provenance without executing work or granting authority.
