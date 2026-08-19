# Contributing

Open an issue before a material behavior or contract change. Pull requests must
reference that issue, remain focused, and include meaningful test evidence.

Before pushing, run:

```shell
go test ./...
go vet ./...
go run ./cmd/gh-pr-quality-gate validate --run-local
python scripts/validate_package.py
```

Maintainers should additionally run the current official Skill and plugin
validators before release.

Do not merge, publish a release, change Draft state, or modify repository rules
without explicit maintainer authorization.
