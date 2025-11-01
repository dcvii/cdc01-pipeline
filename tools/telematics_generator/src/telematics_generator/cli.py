from __future__ import annotations

import json
import sys
import time
from pathlib import Path
from typing import List, Optional

import typer

from .emitter import TelematicsEmitter


def _resolve_output(path: str) -> Optional[Path]:
    if path == "-":
        return None
    output = Path(path)
    output.parent.mkdir(parents=True, exist_ok=True)
    return output


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
) -> None:
    """Continuously stream JSON Lines telematics events."""

    emitter = TelematicsEmitter(seed=seed, fleet=vin)
    output_path = _resolve_output(out)

    def write_events(events: List[dict]) -> None:
        payload = "\n".join(json.dumps(evt, sort_keys=True) for evt in events) + "\n"
        if output_path is None:
            sys.stdout.write(payload)
            sys.stdout.flush()
        else:
            with output_path.open("a", encoding="utf-8") as fh:
                fh.write(payload)

    batch_count = 0
    while True:
        events = list(emitter.events(batch_size=batch_size))
        write_events(events)
        batch_count += 1

        if iterations and batch_count >= iterations:
            break
        if interval > 0:
            time.sleep(interval)


def console_main() -> None:
    typer.run(main)


if __name__ == "__main__":
    console_main()
