from __future__ import annotations

from pathlib import Path

from telematics_generator.cli import main


def test_cli_writes_rotating_files(tmp_path, monkeypatch) -> None:
    monkeypatch.setattr("telematics_generator.cli.time.sleep", lambda _: None)
    output = tmp_path / "events.jsonl"

    main(
        interval=0.0,
        out=str(output),
        batch_size=2,
        iterations=3,
        seed=101,
        max_file_events=2,
        validate=True,
        vin=None,
        center_lat=37.7749,
        center_lon=-122.4194,
    )

    files = sorted(tmp_path.glob("*.jsonl"))
    assert len(files) == 3

    total_lines = 0
    for path in files:
        with path.open("r", encoding="utf-8") as fh:
            total_lines += sum(1 for _ in fh)

    assert total_lines == 6
