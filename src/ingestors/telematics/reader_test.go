package telematics

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReaderSingleFile(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.jsonl")

	testData := `{"event_id":"test-1","vin":"1HGCM82633A123456","trip_id":"trip-1","event_ts":"2024-03-18T14:07:26Z","sample_gap_sec":5,"latitude":37.7749,"longitude":-122.4194,"altitude_m":23.4,"speed_kph":72.3,"heading_deg":184.0,"odometer_km":15482.6,"engine_rpm":2450,"throttle_pct":38.5,"brake_active":false,"fuel_level_pct":62.4,"coolant_temp_c":89.2,"battery_voltage_v":13.8,"tire_pressure_kpa":{"fl":220,"fr":218,"rl":224,"rr":221},"accel_m_s2":{"x":-0.12,"y":0.03,"z":9.77},"geo_hash":"9q8yyj0"}
{"event_id":"test-2","vin":"5XYZT3LB1LG123789","trip_id":"trip-2","event_ts":"2024-03-18T14:07:31Z","sample_gap_sec":5,"latitude":40.7128,"longitude":-74.0060,"altitude_m":10.5,"speed_kph":45.2,"heading_deg":90.0,"odometer_km":8921.3,"engine_rpm":1800,"throttle_pct":25.0,"brake_active":false,"fuel_level_pct":78.5,"coolant_temp_c":85.0,"battery_voltage_v":14.1,"tire_pressure_kpa":{"fl":215,"fr":215,"rl":220,"rr":220},"accel_m_s2":{"x":0.05,"y":-0.01,"z":9.81},"geo_hash":"dr5regw"}
`
	if err := os.WriteFile(testFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create reader
	reader, err := NewReader(ReaderConfig{
		InputPath: testFile,
		BatchSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}

	// Read events
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventCh, errCh := reader.ReadEvents(ctx)

	totalEvents := 0
	for batch := range eventCh {
		totalEvents += len(batch)

		// Validate first event
		if len(batch) > 0 {
			event := batch[0]
			if event.EventID == "" {
				t.Error("Event ID should not be empty")
			}
			if event.VIN == "" {
				t.Error("VIN should not be empty")
			}
		}
	}

	// Check for errors
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Reader error: %v", err)
		}
	default:
	}

	if totalEvents != 2 {
		t.Errorf("Expected 2 events, got %d", totalEvents)
	}
}

func TestReaderDirectory(t *testing.T) {
	// Create temporary directory with multiple files
	tmpDir := t.TempDir()

	files := []struct {
		name    string
		content string
	}{
		{
			"file1.jsonl",
			`{"event_id":"test-1","vin":"1HGCM82633A123456","trip_id":"trip-1","event_ts":"2024-03-18T14:07:26Z","sample_gap_sec":5,"latitude":37.7749,"longitude":-122.4194,"altitude_m":23.4,"speed_kph":72.3,"heading_deg":184.0,"odometer_km":15482.6,"engine_rpm":2450,"throttle_pct":38.5,"brake_active":false,"fuel_level_pct":62.4,"coolant_temp_c":89.2,"battery_voltage_v":13.8,"tire_pressure_kpa":{"fl":220,"fr":218,"rl":224,"rr":221},"accel_m_s2":{"x":-0.12,"y":0.03,"z":9.77},"geo_hash":"9q8yyj0"}`,
		},
		{
			"file2.jsonl",
			`{"event_id":"test-2","vin":"5XYZT3LB1LG123789","trip_id":"trip-2","event_ts":"2024-03-18T14:07:31Z","sample_gap_sec":5,"latitude":40.7128,"longitude":-74.0060,"altitude_m":10.5,"speed_kph":45.2,"heading_deg":90.0,"odometer_km":8921.3,"engine_rpm":1800,"throttle_pct":25.0,"brake_active":false,"fuel_level_pct":78.5,"coolant_temp_c":85.0,"battery_voltage_v":14.1,"tire_pressure_kpa":{"fl":215,"fr":215,"rl":220,"rr":220},"accel_m_s2":{"x":0.05,"y":-0.01,"z":9.81},"geo_hash":"dr5regw"}`,
		},
	}

	for _, f := range files {
		path := filepath.Join(tmpDir, f.name)
		if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", f.name, err)
		}
	}

	// Create reader
	reader, err := NewReader(ReaderConfig{
		InputPath: tmpDir,
		BatchSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}

	// Read events
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventCh, errCh := reader.ReadEvents(ctx)

	totalEvents := 0
	for batch := range eventCh {
		totalEvents += len(batch)
	}

	// Check for errors
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Reader error: %v", err)
		}
	default:
	}

	if totalEvents != 2 {
		t.Errorf("Expected 2 events from 2 files, got %d", totalEvents)
	}
}

func TestReaderBatching(t *testing.T) {
	// Create temporary test file with 5 events
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.jsonl")

	testData := `{"event_id":"1","vin":"VIN1","trip_id":"trip-1","event_ts":"2024-03-18T14:07:26Z","sample_gap_sec":5,"latitude":37.7749,"longitude":-122.4194,"altitude_m":23.4,"speed_kph":72.3,"heading_deg":184.0,"odometer_km":15482.6,"engine_rpm":2450,"throttle_pct":38.5,"brake_active":false,"fuel_level_pct":62.4,"coolant_temp_c":89.2,"battery_voltage_v":13.8,"tire_pressure_kpa":{"fl":220,"fr":218,"rl":224,"rr":221},"accel_m_s2":{"x":-0.12,"y":0.03,"z":9.77},"geo_hash":"9q8yyj0"}
{"event_id":"2","vin":"VIN2","trip_id":"trip-1","event_ts":"2024-03-18T14:07:26Z","sample_gap_sec":5,"latitude":37.7749,"longitude":-122.4194,"altitude_m":23.4,"speed_kph":72.3,"heading_deg":184.0,"odometer_km":15482.6,"engine_rpm":2450,"throttle_pct":38.5,"brake_active":false,"fuel_level_pct":62.4,"coolant_temp_c":89.2,"battery_voltage_v":13.8,"tire_pressure_kpa":{"fl":220,"fr":218,"rl":224,"rr":221},"accel_m_s2":{"x":-0.12,"y":0.03,"z":9.77},"geo_hash":"9q8yyj0"}
{"event_id":"3","vin":"VIN3","trip_id":"trip-1","event_ts":"2024-03-18T14:07:26Z","sample_gap_sec":5,"latitude":37.7749,"longitude":-122.4194,"altitude_m":23.4,"speed_kph":72.3,"heading_deg":184.0,"odometer_km":15482.6,"engine_rpm":2450,"throttle_pct":38.5,"brake_active":false,"fuel_level_pct":62.4,"coolant_temp_c":89.2,"battery_voltage_v":13.8,"tire_pressure_kpa":{"fl":220,"fr":218,"rl":224,"rr":221},"accel_m_s2":{"x":-0.12,"y":0.03,"z":9.77},"geo_hash":"9q8yyj0"}
{"event_id":"4","vin":"VIN4","trip_id":"trip-1","event_ts":"2024-03-18T14:07:26Z","sample_gap_sec":5,"latitude":37.7749,"longitude":-122.4194,"altitude_m":23.4,"speed_kph":72.3,"heading_deg":184.0,"odometer_km":15482.6,"engine_rpm":2450,"throttle_pct":38.5,"brake_active":false,"fuel_level_pct":62.4,"coolant_temp_c":89.2,"battery_voltage_v":13.8,"tire_pressure_kpa":{"fl":220,"fr":218,"rl":224,"rr":221},"accel_m_s2":{"x":-0.12,"y":0.03,"z":9.77},"geo_hash":"9q8yyj0"}
{"event_id":"5","vin":"VIN5","trip_id":"trip-1","event_ts":"2024-03-18T14:07:26Z","sample_gap_sec":5,"latitude":37.7749,"longitude":-122.4194,"altitude_m":23.4,"speed_kph":72.3,"heading_deg":184.0,"odometer_km":15482.6,"engine_rpm":2450,"throttle_pct":38.5,"brake_active":false,"fuel_level_pct":62.4,"coolant_temp_c":89.2,"battery_voltage_v":13.8,"tire_pressure_kpa":{"fl":220,"fr":218,"rl":224,"rr":221},"accel_m_s2":{"x":-0.12,"y":0.03,"z":9.77},"geo_hash":"9q8yyj0"}
`
	if err := os.WriteFile(testFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create reader with batch size of 2
	reader, err := NewReader(ReaderConfig{
		InputPath: testFile,
		BatchSize: 2,
	})
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}

	// Read events
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventCh, errCh := reader.ReadEvents(ctx)

	batches := 0
	totalEvents := 0
	for batch := range eventCh {
		batches++
		totalEvents += len(batch)

		// First two batches should have 2 events, last batch should have 1
		if batches < 3 && len(batch) != 2 {
			t.Errorf("Batch %d: expected 2 events, got %d", batches, len(batch))
		}
		if batches == 3 && len(batch) != 1 {
			t.Errorf("Final batch: expected 1 event, got %d", len(batch))
		}
	}

	// Check for errors
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Reader error: %v", err)
		}
	default:
	}

	if batches != 3 {
		t.Errorf("Expected 3 batches, got %d", batches)
	}
	if totalEvents != 5 {
		t.Errorf("Expected 5 total events, got %d", totalEvents)
	}
}

func TestReaderEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.jsonl")

	if err := os.WriteFile(testFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	reader, err := NewReader(ReaderConfig{
		InputPath: testFile,
		BatchSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventCh, errCh := reader.ReadEvents(ctx)

	totalEvents := 0
	for batch := range eventCh {
		totalEvents += len(batch)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Reader error: %v", err)
		}
	default:
	}

	if totalEvents != 0 {
		t.Errorf("Expected 0 events from empty file, got %d", totalEvents)
	}
}

func TestReaderInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "invalid.jsonl")

	testData := `{"event_id":"test-1","vin":"123"}
invalid json line
{"event_id":"test-2","vin":"456"}
`
	if err := os.WriteFile(testFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	reader, err := NewReader(ReaderConfig{
		InputPath: testFile,
		BatchSize: 10,
	})
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventCh, errCh := reader.ReadEvents(ctx)

	// Drain event channel
	for range eventCh {
	}

	// Should get an error due to invalid JSON
	select {
	case err := <-errCh:
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	case <-time.After(1 * time.Second):
		t.Error("Expected error but got timeout")
	}
}
