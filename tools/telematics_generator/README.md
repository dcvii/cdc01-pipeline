# Telematics Generator

Run synthetic auto telematics events using `uv`:

```sh
uv run --project tools/telematics_generator telematics-generator \
  --interval 1 \
  --out data/raw/telematics/events.jsonl \
  --batch-size 5 \
  --max-file-events 1000
```

Helpful flags:
- `--vin` to pin a fleet (`--vin 1HG... --vin 5XT...`).
- `--center-lat` / `--center-lon` to relocate the fleet.
- `--seed` and `--batch-size` for deterministic replay and volume tuning.
- `--no-validate` if you want to skip schema checks for max throughput.
