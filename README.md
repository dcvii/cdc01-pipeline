# Temporal CDC Pipeline

A Temporal-orchestrated streaming data platform for ingesting, processing, and managing telematics data using Apache Iceberg tables.

## Architecture

This project implements a Bronze → Silver → Gold data architecture:
- **Bronze**: Raw event logs stored as Iceberg tables in object storage (S3/MinIO)
- **Silver**: Conformed, deduplicated data with MERGE/upsert handling
- **Gold/Hot**: Selective minimal subsets in query-optimized stores (Databricks/DuckDB/Postgres)

See [docs/strategy.md](docs/strategy.md) for the complete architectural overview.

## Project Structure

```
.
├── cmd/                    # Command-line applications
│   └── telematics-ingestor/   # Telematics data ingestor CLI
├── context/                # Project documentation and planning
│   ├── AGENTS.md           # Repository guidelines
│   ├── data_source.md      # Telematics schema specification
│   └── plan.md             # Project roadmap
├── docs/                   # Additional documentation
├── pkg/                    # Shared libraries and utilities
│   └── models/             # Data models (telematics events)
├── src/
│   ├── ingestors/          # Source-specific ingestors
│   │   └── telematics/     # Telematics event ingestor
│   └── workflows/          # Temporal workflows and activities
├── tools/                  # Development tools
│   └── telematics_generator/  # Python-based test data generator
└── bin/                    # Compiled binaries
```

## Getting Started

### Prerequisites

- Go 1.22+
- Python 3.12+ with `uv` (for data generation)
- Docker Compose (for Temporal, Redpanda, MinIO)

### Build

```bash
# Build the telematics ingestor
go build -o bin/telematics-ingestor ./cmd/telematics-ingestor

# Run tests
go test ./...

# Lint and format
golangci-lint run
gofmt -w .
```

## Telematics Ingestor

The telematics ingestor reads JSON Lines files containing vehicle telematics events and processes them through configurable output handlers.

### Usage

```bash
# Read from a single file
./bin/telematics-ingestor --input data/raw/telematics/events.jsonl

# Read from a directory (processes all .jsonl files)
./bin/telematics-ingestor --input data/raw/telematics/ --batch-size 50

# Write output to a file
./bin/telematics-ingestor --input events.jsonl --output processed.jsonl

# Verbose mode (prints all events to stdout)
./bin/telematics-ingestor --input events.jsonl --verbose
```

### Options

- `--input` (required): Path to JSON Lines file or directory
- `--output`: Path to output file (defaults to stdout if not specified)
- `--batch-size`: Number of events per batch (default: 100)
- `--verbose`: Enable verbose output (prints all events)

### Example Output

```
2025/10/31 17:53:52 Starting ingestor (input=data/raw/telematics/test_events.jsonl, batch_size=2)
2025/10/31 17:53:52 Processed batch: 2 events (first VIN: 1HGCM82633A123456, last VIN: 1HGCM82633A123456)
2025/10/31 17:53:52 Processed batch: 2 events (first VIN: 5XYZT3LB1LG123789, last VIN: 5XYZT3LB1LG123789)
2025/10/31 17:53:52 Ingestion complete:
2025/10/31 17:53:52   Events processed: 5
2025/10/31 17:53:52   Batches processed: 3
2025/10/31 17:53:52   Errors: 0
```

## Telematics Schema

The ingestor processes events according to the canonical schema defined in [context/data_source.md](context/data_source.md).

Key fields:
- **Identifiers**: `event_id` (UUID), `vin` (Vehicle ID), `trip_id` (session)
- **Location**: `latitude`, `longitude`, `altitude_m`, `geo_hash`
- **Motion**: `speed_kph`, `heading_deg`, `odometer_km`
- **Engine**: `engine_rpm`, `throttle_pct`, `fuel_level_pct`, `coolant_temp_c`
- **Sensors**: `tire_pressure_kpa` (per-wheel), `accel_m_s2` (3-axis)

## Data Generation

Generate test telematics data using the Python generator:

```bash
uv run --project tools/telematics_generator telematics-generator \
  --interval 1 \
  --out data/raw/telematics/events.jsonl \
  --batch-size 5 \
  --iterations 10 \
  --seed 42
```

See [docs/data-generator.md](docs/data-generator.md) for details.

## Development Workflow

1. **Create a feature branch**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make changes and test**
   ```bash
   go test ./...
   go build -o bin/telematics-ingestor ./cmd/telematics-ingestor
   ```

3. **Lint and format**
   ```bash
   golangci-lint run
   gofmt -w .
   ```

4. **Commit using Conventional Commits**
   ```bash
   git commit -m "feat: add new ingestor feature"
   ```

## Next Steps

- [ ] Integrate with Kafka/Redpanda for event streaming
- [ ] Implement Iceberg Bronze table writer
- [ ] Create Temporal workflows for ingestion orchestration
- [ ] Add Silver layer transformation logic
- [ ] Implement data quality checks and metrics

## Contributing

See [context/AGENTS.md](context/AGENTS.md) for coding standards, testing guidelines, and contribution workflow.

## License

Internal project - All rights reserved
