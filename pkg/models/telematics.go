package models

import "time"

// TirePressure represents tire pressure readings for all four wheels in kilopascals.
type TirePressure struct {
	FL float64 `json:"fl"` // Front Left
	FR float64 `json:"fr"` // Front Right
	RL float64 `json:"rl"` // Rear Left
	RR float64 `json:"rr"` // Rear Right
}

// Acceleration represents 3-axis acceleration in m/s² relative to vehicle frame.
type Acceleration struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// TelematicsEvent represents a single telematics event from a vehicle.
// Schema defined in context/data_source.md
type TelematicsEvent struct {
	// Core identifiers
	EventID   string    `json:"event_id"`   // UUID for idempotent ingestion
	VIN       string    `json:"vin"`        // Vehicle Identification Number
	TripID    string    `json:"trip_id"`    // Logical trip/session identifier
	EventTS   time.Time `json:"event_ts"`   // Event timestamp (ISO-8601 UTC)
	SampleGap int       `json:"sample_gap_sec"` // Seconds since previous event

	// Location
	Latitude  float64 `json:"latitude"`   // Decimal degrees (-90 to 90)
	Longitude float64 `json:"longitude"`  // Decimal degrees (-180 to 180)
	AltitudeM float64 `json:"altitude_m"` // Meters above sea level (optional)
	GeoHash   string  `json:"geo_hash"`   // Precision-7 geohash for partitioning

	// Motion
	SpeedKPH   float64 `json:"speed_kph"`   // Kilometers per hour
	HeadingDeg float64 `json:"heading_deg"` // Compass heading (0-360)
	OdometerKM float64 `json:"odometer_km"` // Cumulative odometer reading

	// Engine/Drivetrain
	EngineRPM   int     `json:"engine_rpm"`   // Revolutions per minute
	ThrottlePct float64 `json:"throttle_pct"` // Accelerator position (0-100)
	BrakeActive bool    `json:"brake_active"` // Brake pedal depressed

	// Fuel and Temperature
	FuelLevelPct  float64 `json:"fuel_level_pct"`  // Fuel remaining (0-100)
	CoolantTempC  float64 `json:"coolant_temp_c"`  // Engine coolant temp (Celsius)
	BatteryVoltV  float64 `json:"battery_voltage_v"` // 12V system voltage

	// Sensors
	TirePressureKPA TirePressure  `json:"tire_pressure_kpa"` // Per-wheel tire pressure
	AccelMS2        Acceleration  `json:"accel_m_s2"`        // 3-axis acceleration
}
