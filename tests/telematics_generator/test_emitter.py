from __future__ import annotations

from telematics_generator.emitter import TelematicsEmitter
from telematics_generator.schema import SCHEMA


def test_emitter_respects_schema_keys() -> None:
    emitter = TelematicsEmitter(seed=42, fleet=["TESTVIN1234567890"])
    event = next(iter(emitter.events(batch_size=1)))

    missing = set(SCHEMA.keys()) - event.keys()
    assert not missing, f"Emitter missing keys: {missing}"


def test_emitter_repeatable_with_seed() -> None:
    emitter_a = TelematicsEmitter(seed=123, fleet=["VIN123"])
    emitter_b = TelematicsEmitter(seed=123, fleet=["VIN123"])

    events_a = [next(iter(emitter_a.events(batch_size=1))) for _ in range(3)]
    events_b = [next(iter(emitter_b.events(batch_size=1))) for _ in range(3)]

    assert events_a == events_b
