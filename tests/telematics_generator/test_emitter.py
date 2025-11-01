from __future__ import annotations

from typing import List

from telematics_generator.emitter import TelematicsEmitter
from telematics_generator.schema import SCHEMA, validate_event


def test_emitter_respects_schema_keys() -> None:
    emitter = TelematicsEmitter(seed=42, fleet=["TESTVIN1234567890"])
    event = next(emitter.events(batch_size=1))

    missing = set(SCHEMA.keys()) - event.keys()
    assert not missing, f"Emitter missing keys: {missing}"
    validate_event(event)


def test_emitter_repeatable_with_seed() -> None:
    emitter_a = TelematicsEmitter(seed=123, fleet=["VIN123"])
    emitter_b = TelematicsEmitter(seed=123, fleet=["VIN123"])

    events_a = [next(emitter_a.events(batch_size=1)) for _ in range(3)]
    events_b = [next(emitter_b.events(batch_size=1)) for _ in range(3)]

    comparable_keys = set(events_a[0].keys()) - {"event_id", "trip_id", "event_ts"}
    trimmed_a = [{key: event[key] for key in comparable_keys} for event in events_a]
    trimmed_b = [{key: event[key] for key in comparable_keys} for event in events_b]

    assert trimmed_a == trimmed_b


def test_emitter_batch_size_exceeded() -> None:
    fleet: List[str] = [f"VIN{i:02d}VIN000000" for i in range(2)]
    emitter = TelematicsEmitter(seed=999, fleet=fleet)

    events = list(emitter.events(batch_size=5))
    assert len(events) == 5
    vins = {event["vin"] for event in events}
    assert vins <= set(fleet)
