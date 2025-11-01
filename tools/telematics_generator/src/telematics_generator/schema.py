"""Schema definition and helper utilities for telematics events."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any, Dict, Mapping


@dataclass(frozen=True)
class FieldSpec:
    """Metadata for a schema field."""

    name: str
    type_hint: str
    description: str
    example: Any


SCHEMA: Mapping[str, FieldSpec] = {
    "event_id": FieldSpec(
        name="event_id",
        type_hint="str (uuid)",
        description="Unique identifier per emitted event.",
        example="5c7a9b84-6d1f-4b0d-91a8-b09e4395f822",
    ),
    "vin": FieldSpec(
        name="vin",
        type_hint="str (17)",
        description="Vehicle identification number.",
        example="1HGCM82633A123456",
    ),
    "trip_id": FieldSpec(
        name="trip_id",
        type_hint="str (uuid)",
        description="Logical trip/session identifier.",
        example="bd9d3eda-02af-40fd-b5fd-9fef77e1c4da",
    ),
    "event_ts": FieldSpec(
        name="event_ts",
        type_hint="str (iso-8601)",
        description="Event timestamp in UTC.",
        example="2024-03-18T14:07:26Z",
    ),
    "sample_gap_sec": FieldSpec(
        name="sample_gap_sec",
        type_hint="int",
        description="Seconds since last event for VIN.",
        example=5,
    ),
    "latitude": FieldSpec(
        name="latitude",
        type_hint="float",
        description="Latitude in decimal degrees.",
        example=37.7749,
    ),
    "longitude": FieldSpec(
        name="longitude",
        type_hint="float",
        description="Longitude in decimal degrees.",
        example=-122.4194,
    ),
    "altitude_m": FieldSpec(
        name="altitude_m",
        type_hint="float",
        description="Altitude above sea level in meters.",
        example=23.4,
    ),
    "speed_kph": FieldSpec(
        name="speed_kph",
        type_hint="float",
        description="Vehicle speed in kilometers per hour.",
        example=72.3,
    ),
    "heading_deg": FieldSpec(
        name="heading_deg",
        type_hint="float",
        description="Compass heading in degrees.",
        example=184.0,
    ),
    "odometer_km": FieldSpec(
        name="odometer_km",
        type_hint="float",
        description="Cumulative odometer reading in kilometers.",
        example=15482.6,
    ),
    "engine_rpm": FieldSpec(
        name="engine_rpm",
        type_hint="int",
        description="Engine speed in RPM.",
        example=2450,
    ),
    "throttle_pct": FieldSpec(
        name="throttle_pct",
        type_hint="float",
        description="Throttle pedal position in percent.",
        example=38.5,
    ),
    "brake_active": FieldSpec(
        name="brake_active",
        type_hint="bool",
        description="Brake pedal depressed indicator.",
        example=False,
    ),
    "fuel_level_pct": FieldSpec(
        name="fuel_level_pct",
        type_hint="float",
        description="Fuel tank level percentage.",
        example=62.4,
    ),
    "coolant_temp_c": FieldSpec(
        name="coolant_temp_c",
        type_hint="float",
        description="Engine coolant temperature in Celsius.",
        example=89.2,
    ),
    "battery_voltage_v": FieldSpec(
        name="battery_voltage_v",
        type_hint="float",
        description="Battery voltage in volts.",
        example=13.8,
    ),
    "tire_pressure_kpa": FieldSpec(
        name="tire_pressure_kpa",
        type_hint="object",
        description="Per-wheel tire pressure in kilopascals.",
        example={"fl": 220, "fr": 218, "rl": 224, "rr": 221},
    ),
    "accel_m_s2": FieldSpec(
        name="accel_m_s2",
        type_hint="object",
        description="3-axis acceleration in m/s^2.",
        example={"x": -0.12, "y": 0.03, "z": 9.77},
    ),
    "geo_hash": FieldSpec(
        name="geo_hash",
        type_hint="str (7)",
        description="Geohash derived from lat/lon.",
        example="9q8yyj0",
    ),
}


def schema_examples() -> Dict[str, Any]:
    """Return a mapping of example values keyed by field name."""

    return {field: spec.example for field, spec in SCHEMA.items()}
