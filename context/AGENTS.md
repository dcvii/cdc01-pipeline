# Repository Guidelines

## Project Structure & Module Organization
Keep architecture notes in `strategy.md`; use it as the contract for new code. Place Temporal workers and shared workflow logic under `src/workflows/`, source-specific ingestors in `src/ingestors/`, and reusable utilities (Iceberg helpers, schema clients) in `pkg/`. Infrastructure files (`docker-compose.yaml`, Terraform, Helm) belong in `infra/`. Keep data-contract specs (`.proto`, Avro) in `contracts/` paired with generated code in `gen/`. Documentation and runbooks should live in `docs/`.

## Build, Test, and Development Commands
Run Temporal workers locally with:
```sh
go run ./cmd/workers
```
Execute unit tests:
```sh
go test ./...
```
Lint and format before pushing:
```sh
golangci-lint run && gofmt -w .
```
Spin up dependencies (Temporal, Redpanda, MinIO) via:
```sh
docker compose up temporal redpanda minio
```

## Coding Style & Naming Conventions
Use Go 1.22+ tooling. Keep imports grouped and sorted (`goimports`). Stick to `gofmt` for indentation (tabs). Name workflows with the suffix `Workflow`, activities with `Activity`, and commands/config structs with descriptive nouns (`CdcStreamConfig`). Configuration files should use snake_case (e.g., `cdc_pipeline.yaml`), whereas Go packages stay lower_snake or single words.

## Testing Guidelines
Write table-driven Go tests named `Test<Entity><Behavior>`. Cover happy paths, idempotency, and failure retries for every workflow/activity. Integration suites that touch Temporal or Redpanda should live under `test/integration/` and use Docker containers spun up by `docker compose`. Track coverage; keep unit coverage above 80% for workflow logic, and document any gaps in `docs/test-notes.md`.

## Commit & Pull Request Guidelines
Follow Conventional Commits (`feat:`, `fix:`, `chore:`, `docs:`). Reference Temporal workflow IDs or schema versions when relevant (`feat: add bronze compaction workflow`). In pull requests, include: purpose summary, linked issue, validation commands (`go test`, `docker compose ps`), and screenshots/log excerpts for pipeline changes. Request review from data platform and infra maintainers when touching `infra/` or `contracts/`.

## Temporal & Data Pipeline Notes
Schema evolution demands backward compatibility; update contracts first, regenerate code, then roll workers. Treat Iceberg tables as the source of truth—avoid writing directly to hot stores outside orchestrated activities. Use Temporal search attributes for observability (`source`, `watermark`) and document any new metrics or alerts in `docs/observability.md`.
