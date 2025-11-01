"""Generate synthetic telematics events that conform to the documented schema."""

from __future__ import annotations

import random
import uuid
from dataclasses import dataclass
from typing import Dict, Iterable, Iterator, List, Optional

import geohash2
import pendulum
from faker import Faker


@dataclass
class VehicleState:
    """Track per-vehicle state to keep values realistic."""

    vin: str
    trip_id: str
    last_odometer_km: float
    last_ts: pendulum.DateTime
    home_lat: float
    home_lon: float


class TelematicsEmitter:
    """Create synthetic telematics events with light state management."""

    def __init__(
        self,
        seed: Optional[int] = None,
        fleet: Optional[Iterable[str]] = None,
        base_location: tuple[float, float] = (37.7749, -122.4194),
    ) -> None:
        self._faker = Faker()
        self._rng = random.Random(seed)
        if seed is not None:
            self._faker.seed_instance(seed)
        self._base_lat, self._base_lon = base_location

        if fleet:
            self._fleet: List[str] = list(fleet)
        else:
            # Default fleet of 5 VINs if none provided.
            self._fleet = [self._faker.unique.vin() for _ in range(5)]

        now = pendulum.now("UTC")
        self._state: Dict[str, VehicleState] = {
            vin: VehicleState(
                vin=vin,
                trip_id=str(uuid.uuid4()),
                last_odometer_km=self._rng.uniform(0, 200_000) / 10.0,
                last_ts=now.subtract(minutes=self._rng.randint(0, 15)),
                home_lat=self._base_lat + self._rng.uniform(-0.05, 0.05),
                home_lon=self._base_lon + self._rng.uniform(-0.05, 0.05),
            )
            for vin in self._fleet
        }

    def events(self, batch_size: int = 1) -> Iterator[dict]:
        """Yield a batch of telematics events across the fleet."""

        if not self._fleet:
            return

        for vin in self._rng.choices(self._fleet, k=batch_size):
            yield self._next_event(vin)

    def _next_event(self, vin: str) -> dict:
        state = self._state[vin]

        gap_seconds = self._rng.choice((1, 2, 3, 5, 10))
        ts = state.last_ts.add(seconds=gap_seconds)

        # Movement deltas
        speed_kph = max(0.0, self._rng.gauss(65, 20))
        heading_deg = self._rng.uniform(0, 360)
        odometer = max(state.last_odometer_km + speed_kph * (gap_seconds / 3600.0), 0.0)

        latitude = state.home_lat + self._rng.uniform(-0.01, 0.01)
        longitude = state.home_lon + self._rng.uniform(-0.01, 0.01)
        altitude_m = max(0.0, self._rng.gauss(20, 5))

        tire_pressure = {
            "fl": self._rng.randint(200, 240),
            "fr": self._rng.randint(200, 240),
            "rl": self._rng.randint(200, 240),
            "rr": self._rng.randint(200, 240),
        }

        accel = {axis: round(self._rng.uniform(-1.5, 1.5), 2) for axis in ("x", "y", "z")}
        accel["z"] = round(9.81 + self._rng.uniform(-0.2, 0.2), 2)

        event = {
            "event_id": str(uuid.uuid4()),
            "vin": vin,
            "trip_id": state.trip_id,
            "event_ts": ts.to_iso8601_string(),
            "sample_gap_sec": gap_seconds,
            "latitude": round(latitude, 6),
            "longitude": round(longitude, 6),
            "altitude_m": round(altitude_m, 1),
            "speed_kph": round(speed_kph, 1),
            "heading_deg": round(heading_deg, 1),
            "odometer_km": round(odometer, 1),
            "engine_rpm": int(max(0, self._rng.gauss(2200, 600))),
            "throttle_pct": round(self._rng.uniform(0, 100), 1),
            "brake_active": self._rng.random() < 0.1,
            "fuel_level_pct": round(max(0.0, self._rng.gauss(60, 5)), 1),
            "coolant_temp_c": round(self._rng.gauss(90, 5), 1),
            "battery_voltage_v": round(self._rng.gauss(13.8, 0.2), 2),
            "tire_pressure_kpa": tire_pressure,
            "accel_m_s2": accel,
            "geo_hash": geohash2.encode(latitude, longitude, precision=7),
        }

        # Update vehicle state
        if self._rng.random() < 0.02:
            # Occasionally roll a new trip id to simulate ignition cycle.
            state.trip_id = str(uuid.uuid4())

        state.last_ts = ts
        state.last_odometer_km = odometer
        self._state[vin] = state
        return event
