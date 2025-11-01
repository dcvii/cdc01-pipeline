# Auto Telematics Schema

Field | Type | Example | Notes
--- | --- | --- | ---
`event_id` | UUID string | `5c7a9b84-6d1f-4b0d-91a8-b09e4395f822` | Unique identifier for each emission; enables idempotent ingestion.
`vin` | String (17) | `1HGCM82633A123456` | Uppercase alphanumeric Vehicle Identification Number.
`trip_id` | UUID string | `bd9d3eda-02af-40fd-b5fd-9fef77e1c4da` | Logical trip/session identifier; resets when ignition cycles.
`event_ts` | ISO-8601 UTC string | `2024-03-18T14:07:26Z` | Event timestamp truncated to milliseconds.
`sample_gap_sec` | Integer | `5` | Seconds since previous event for this VIN.
`latitude` | Float | `37.7749` | Decimal degrees; valid range -90 to 90.
`longitude` | Float | `-122.4194` | Decimal degrees; valid range -180 to 180.
`altitude_m` | Float | `23.4` | Optional; meters above sea level.
`speed_kph` | Float | `72.3` | Kilometers per hour; 0-250 typical for passenger vehicles.
`heading_deg` | Float | `184.0` | Compass heading in degrees (0-360).
`odometer_km` | Float | `15482.6` | Cumulative odometer reading in kilometers.
`engine_rpm` | Integer | `2450` | Revolutions per minute; 0-8000 typical.
`throttle_pct` | Float | `38.5` | Accelerator position percentage (0-100).
`brake_active` | Boolean | `false` | `true` when brake pedal depressed.
`fuel_level_pct` | Float | `62.4` | Percent of fuel remaining (0-100).
`coolant_temp_c` | Float | `89.2` | Engine coolant temperature in Celsius.
`battery_voltage_v` | Float | `13.8` | 12V system voltage.
`tire_pressure_kpa` | Object | `{"fl": 220, "fr": 218, "rl": 224, "rr": 221}` | Tire pressure per wheel in kilopascals.
`accel_m_s2` | Object | `{"x": -0.12, "y": 0.03, "z": 9.77}` | 3-axis acceleration relative to vehicle frame.
`geo_hash` | String (7) | `9q8yyj0` | Precision-7 geohash derived from lat/lon for partitioning.

All numeric values should include realistic jitter suited for downstream analytics. Objects serialize as nested JSON values in the emitted record.
