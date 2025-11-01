# Telematics Generator

Run synthetic auto telematics events using `uv`:

```sh
uv run telematics-generator --interval 5 --out data/raw/telematics
```

Use `--vin` to pin specific vehicles, `--seed` for deterministic runs, and `--batch-size` to emit multiple events per interval.
