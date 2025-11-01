from __future__ import annotations

import json
import sys
import time
from pathlib import Path
from typing import List, Optional

import typer

from .emitter import TelematicsEmitter
from .schema import validate_event

__all__ = ["main", "console_main"]


def _resolve_output(path: str) -> Optional[Path]:
    if path == "-":
        return None
    output = Path(path)
    output.parent.mkdir(parents=True, exist_ok=True)
    return output


class _RotatingWriter:
    """Handle optional JSONL file rotation."""

    def __init__(self, base_path: Path, max_events: int) -> None:
        self._base_path = base_path
        self._max_events = max_events if max_events > 0 else 0
        self._current_path = base_path
        self._events_written = 0
        self._file_index = 0
        self._handle = None

    def _ensure_handle(self) -> None:
        if self._handle is None:
            self._handle = self._current_path.open("a", encoding="utf-8")

    def _next_path(self) -> Path:
        suffix = self._base_path.suffix or ".jsonl"
        stem = self._base_path.stem
        self._file_index += 1
        return self._base_path.with_name(f"{stem}_{self._file_index:04d}{suffix}")

    def write_line(self, line: str) -> None:
        if self._max_events and self._events_written >= self._max_events:
            self._close()
            self._current_path = self._next_path()
            self._events_written = 0

        self._ensure_handle()
        assert self._handle is not None  # for type checkers
        self._handle.write(line)
        self._handle.flush()
        self._events_written += 1

    def _close(self) -> None:
        if self._handle:
            self._handle.close()
            self._handle = None

    def close(self) -> None:
        self._close()


def main(  # noqa: PLR0913 - CLI surface area is intentional
    interval: float = typer.Option(
        5.0, help="Seconds to sleep between batches.", min=0.0
    ),
    out: str = typer.Option("-", help="Output file (use - for stdout)."),
    batch_size: int = typer.Option(
        1, min=1, max=1000, help="Number of VIN events to emit per batch."
    ),
    iterations: int = typer.Option(
        0, min=0, help="Number of batches to emit (0 for infinite)."
    ),
    seed: Optional[int] = typer.Option(
        None, help="Random seed for deterministic output."
    ),
    vin: Optional[List[str]] = typer.Option(
        None,
        "--vin",
        "-v",
        help="Explicit VINs to simulate. Repeat flag for multiple vehicles.",
    ),
    center_lat: float = typer.Option(
        37.7749, help="Base latitude used to center generated coordinates."
    ),
    center_lon: float = typer.Option(
        -122.4194, help="Base longitude used to center generated coordinates."
    ),
    validate: bool = typer.Option(
        True, "--validate/--no-validate", help="Validate events before emitting."
    ),
    max_file_events: int = typer.Option(
        0,
        help="Rotate file after N events when writing to disk (0 disables rotation).",
    ),
) -> None:
    """Continuously stream JSON Lines telematics events."""

    fleet = [value.upper() for value in vin] if vin else None
    emitter = TelematicsEmitter(
        seed=seed,
        fleet=fleet,
        base_location=(center_lat, center_lon),
    )
    output_path = _resolve_output(out)

    writer: Optional[_RotatingWriter] = None
    if output_path is not None:
        writer = _RotatingWriter(output_path, max_events=max_file_events)

    def write_event(event: dict) -> None:
        if validate:
            try:
                validate_event(event)
            except ValueError as exc:
                if writer:
                    writer.close()
                typer.secho(f"Validation failed: {exc}", err=True, fg=typer.colors.RED)
                raise typer.Exit(code=1) from exc
        payload = json.dumps(event, sort_keys=True) + "\n"
        if output_path is None:
            sys.stdout.write(payload)
            sys.stdout.flush()
        else:
            assert writer is not None
            writer.write_line(payload)

    batch_count = 0
    while True:
        emitted_this_batch = 0
        for event in emitter.events(batch_size=batch_size):
            write_event(event)
            emitted_this_batch += 1

        batch_count += 1

        if iterations and batch_count >= iterations:
            break
        if interval > 0 and emitted_this_batch:
            time.sleep(interval)

    if writer:
        writer.close()


def console_main() -> None:
    typer.run(main)


if __name__ == "__main__":
    console_main()
