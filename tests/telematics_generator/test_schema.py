from __future__ import annotations

import copy

import pytest

from telematics_generator.schema import SCHEMA, schema_examples, validate_event


def test_schema_examples_validate() -> None:
    payload = schema_examples()
    validate_event(payload)


def test_validate_event_missing_field_raises() -> None:
    payload = schema_examples()
    payload.pop("vin")
    with pytest.raises(ValueError):
        validate_event(payload)


def test_validate_event_invalid_tire_pressure() -> None:
    payload = schema_examples()
    payload["tire_pressure_kpa"] = {"fl": "bad"}
    with pytest.raises(ValueError):
        validate_event(payload)
