# Project Plan

## Immediate Focus — Telematics Data Generator
- Finalize the canonical schema in `context/data_source.md` (field names, units, and sample values) so the generator code has a single source of truth.
- Scaffold a Python tool under `tools/telematics_generator/` using `uv` for dependency and script management; pin `faker` and `pendulum` (for time handling) in `pyproject.toml`.
- Implement a CLI (`uv run telematics-generator --interval 5s --out data/raw/telematics`) that streams JSON Lines records on an interval, honoring the schema (VIN, trip_id, timestamp, GPS, speed, heading, telematics metrics).
- Support configurable random seeds and batch sizes to enable repeatable tests; document defaults in `docs/data-generator.md`.
- Add lightweight tests (pytest) to validate field presence, ranges, and deterministic seeding, and automate them with `uv run pytest`.

## Near-Term Follow Ups
- Land generated files into an S3-compatible staging bucket (MinIO locally) and register the Iceberg Bronze table that Temporal will manage.
- Extend schema contracts (`contracts/`) and code generation once the generator fields are stable.
- Instrument generator with basic metrics/logging and surface them through Temporal once ingestion workflows exist.

## Parking Lot
- Temporal workflow to schedule/monitor generator runs alongside other CDC sources.
- Automated promotion from Bronze to Silver tables, including dedupe and enrichment jobs keyed off the generated data.
