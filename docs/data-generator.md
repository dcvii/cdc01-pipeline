# Telematics Data Generator

Run the generator with `uv` from the repo root:

```sh
uv run --project tools/telematics_generator telematics-generator --interval 2 --out data/raw/telematics/telematics.jsonl --iterations 10 --batch-size 3 --seed 42
```

Flags:
- `--interval`: seconds between batches (use `0` to disable sleeping).
- `--out`: `-` for stdout, otherwise path to append JSON Lines output.
- `--batch-size`: number of VIN events per batch.
- `--iterations`: stop after N batches (default `0` emits forever).
- `--vin`: specify VINs to simulate; omit to auto-generate a fleet of five.
- `--seed`: deterministic Faker output.

Schema details live in `context/data_source.md`.
