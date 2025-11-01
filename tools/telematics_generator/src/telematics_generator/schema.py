"""Schema definition, examples, and validation helpers for telematics events."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Callable, Dict, Mapping, Optional

__all__ = ["FieldSpec", "SCHEMA", "schema_examples", "validate_event"]


def _is_numeric(value: Any) -> bool:
    return isinstance(value, (int, float)) and not isinstance(value, bool)


def _is_uuid_str(value: Any) -> bool:
    return isinstance(value, str) and len(value) >= 16


def _is_geohash(value: Any) -> bool:
    return isinstance(value, str) and 1 <= len(value) <= 12


@dataclass(frozen=True)
class FieldSpec:
    """Metadata for a schema field."""

    name: str
    type_hint: str
    description: str
    example: Any
    validator: Optional[Callable[[Any], bool]] = None


SCHEMA: Mapping[str, FieldSpec] = {
    "event_id": FieldSpec(
        name="event_id",
        type_hint="str (uuid)",
        description="Unique identifier per emitted event.",
        example="5c7a9b84-6d1f-4b0d-91a8-b09e4395f822",
        validator=_is_uuid_str,
    ),
    "vin": FieldSpec(
        name="vin",
        type_hint="str (17)",
        description="Vehicle identification number.",
        example="1HGCM82633A123456",
        validator=lambda value: isinstance(value, str) and len(value) == 17,
    ),
    "trip_id": FieldSpec(
        name="trip_id",
        type_hint="str (uuid)",
        description="Logical trip/session identifier.",
        example="bd9d3eda-02af-40fd-b5fd-9fef77e1c4da",
        validator=_is_uuid_str,
    ),
    "event_ts": FieldSpec(
        name="event_ts",
        type_hint="str (iso-8601)",
        description="Event timestamp in UTC.",
        example="2024-03-18T14:07:26Z",
        validator=lambda value: isinstance(value, str) and value.endswith("Z"),
    ),
    "sample_gap_sec": FieldSpec(
        name="sample_gap_sec",
        type_hint="int",
        description="Seconds since last event for VIN.",
        example=5,
        validator=lambda value: isinstance(value, int) and value >= 0,
    ),
    "latitude": FieldSpec(
        name="latitude",
        type_hint="float",
        description="Latitude in decimal degrees.",
        example=37.7749,
        validator=lambda value: _is_numeric(value) and -90 <= float(value) <= 90,
    ),
    "longitude": FieldSpec(
        name="longitude",
        type_hint="float",
        description="Longitude in decimal degrees.",
        example=-122.4194,
        validator=lambda value: _is_numeric(value) and -180 <= float(value) <= 180,
    ),
    "altitude_m": FieldSpec(
        name="altitude_m",
        type_hint="float",
        description="Altitude above sea level in meters.",
        example=23.4,
        validator=_is_numeric,
    ),
    "speed_kph": FieldSpec(
        name="speed_kph",
        type_hint="float",
        description="Vehicle speed in kilometers per hour.",
        example=72.3,
        validator=lambda value: _is_numeric(value) and float(value) >= 0,
    ),
    "heading_deg": FieldSpec(
        name="heading_deg",
        type_hint="float",
        description="Compass heading in degrees.",
        example=184.0,
        validator=_is_numeric,
    ),
    "odometer_km": FieldSpec(
        name="odometer_km",
        type_hint="float",
        description="Cumulative odometer reading in kilometers.",
        example=15482.6,
        validator=lambda value: _is_numeric(value) and float(value) >= 0,
    ),
    "engine_rpm": FieldSpec(
        name="engine_rpm",
        type_hint="int",
        description="Engine speed in RPM.",
        example=2450,
        validator=lambda value: isinstance(value, int) and value >= 0,
    ),
    "throttle_pct": FieldSpec(
        name="throttle_pct",
        type_hint="float",
        description="Throttle pedal position in percent.",
        example=38.5,
        validator=lambda value: _is_numeric(value) and 0 <= float(value) <= 100,
    ),
    "brake_active": FieldSpec(
        name="brake_active",
        type_hint="bool",
        description="Brake pedal depressed indicator.",
        example=False,
        validator=lambda value: isinstance(value, bool),
    ),
    "fuel_level_pct": FieldSpec(
        name="fuel_level_pct",
        type_hint="float",
        description="Fuel tank level percentage.",
        example=62.4,
        validator=lambda value: _is_numeric(value) and 0 <= float(value) <= 100,
    ),
    "coolant_temp_c": FieldSpec(
        name="coolant_temp_c",
        type_hint="float",
        description="Engine coolant temperature in Celsius.",
        example=89.2,
        validator=_is_numeric,
    ),
    "battery_voltage_v": FieldSpec(
        name="battery_voltage_v",
        type_hint="float",
        description="Battery voltage in volts.",
        example=13.8,
        validator=_is_numeric,
    ),
    "tire_pressure_kpa": FieldSpec(
        name="tire_pressure_kpa",
        type_hint="object",
        description="Per-wheel tire pressure in kilopascals.",
        example={"fl": 220, "fr": 218, "rl": 224, "rr": 221},
        validator=lambda value: isinstance(value, Mapping)
        and {"fl", "fr", "rl", "rr"}.issubset(set(value.keys()))
        and all(_is_numeric(value[wheel]) for wheel in ("fl", "fr", "rl", "rr")),
    ),
    "accel_m_s2": FieldSpec(
        name="accel_m_s2",
        type_hint="object",
        description="3-axis acceleration in m/s^2.",
        example={"x": -0.12, "y": 0.03, "z": 9.77},
        validator=lambda value: isinstance(value, Mapping)
        and {"x", "y", "z"}.issubset(set(value.keys()))
        and all(_is_numeric(value[axis]) for axis in ("x", "y", "z")),
    ),
    "geo_hash": FieldSpec(
        name="geo_hash",
        type_hint="str (7)",
        description="Geohash derived from lat/lon.",
        example="9q8yyj0",
        validator=_is_geohash,
    ),
}


def schema_examples() -> Dict[str, Any]:
    """Return a mapping of example values keyed by field name."""

    return {field: spec.example for field, spec in SCHEMA.items()}


def validate_event(payload: Mapping[str, Any]) -> None:
    """Validate that a payload conforms to the declared schema.

    Args:
        payload: Mapping representing an emitted telematics event.

    Raises:
        ValueError: When a required field is missing or fails validation.
    """

    missing = [field for field in SCHEMA if field not in payload]
    if missing:
        raise ValueError(f"Missing schema fields: {missing}")

    for spec in SCHEMA.values():
        value = payload.get(spec.name)
        if spec.validator and not spec.validator(value):
            raise ValueError(
                f"Field '{spec.name}' failed validation "
                f"(value={value!r}, expected={spec.type_hint})"
            )
