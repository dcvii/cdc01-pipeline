package telematics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mdcb/cdc01/pkg/models"
)

// StdoutHandler writes events to stdout (useful for testing and debugging).
type StdoutHandler struct {
	verbose bool
}

// NewStdoutHandler creates an output handler that writes to stdout.
func NewStdoutHandler(verbose bool) *StdoutHandler {
	return &StdoutHandler{verbose: verbose}
}

// ProcessBatch writes a batch of events to stdout.
func (h *StdoutHandler) ProcessBatch(ctx context.Context, events []models.TelematicsEvent) error {
	if h.verbose {
		for _, event := range events {
			data, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal event: %w", err)
			}
			fmt.Println(string(data))
		}
	} else {
		log.Printf("Processed batch: %d events (first VIN: %s, last VIN: %s)",
			len(events),
			events[0].VIN,
			events[len(events)-1].VIN)
	}
	return nil
}

// Close implements OutputHandler.
func (h *StdoutHandler) Close() error {
	return nil
}

// FileHandler writes events to a JSON Lines file.
type FileHandler struct {
	file    *os.File
	encoder *json.Encoder
}

// NewFileHandler creates an output handler that writes to a file.
func NewFileHandler(outputPath string) (*FileHandler, error) {
	f, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	return &FileHandler{
		file:    f,
		encoder: json.NewEncoder(f),
	}, nil
}

// ProcessBatch writes a batch of events to the output file as JSON Lines.
func (h *FileHandler) ProcessBatch(ctx context.Context, events []models.TelematicsEvent) error {
	for _, event := range events {
		if err := h.encoder.Encode(event); err != nil {
			return fmt.Errorf("failed to encode event: %w", err)
		}
	}
	return nil
}

// Close closes the output file.
func (h *FileHandler) Close() error {
	if h.file != nil {
		return h.file.Close()
	}
	return nil
}

// SimpleMetrics is a basic metrics handler that logs statistics.
type SimpleMetrics struct {
	totalEvents  int64
	totalBatches int64
	errors       []error
}

// NewSimpleMetrics creates a new simple metrics handler.
func NewSimpleMetrics() *SimpleMetrics {
	return &SimpleMetrics{
		errors: make([]error, 0),
	}
}

// RecordEventsProcessed records the number of events processed.
func (m *SimpleMetrics) RecordEventsProcessed(count int64) {
	m.totalEvents += count
}

// RecordBatchProcessed records batch processing duration.
func (m *SimpleMetrics) RecordBatchProcessed(durationMs int64) {
	m.totalBatches++
	if durationMs > 1000 {
		log.Printf("Slow batch detected: %dms", durationMs)
	}
}

// RecordError records an error.
func (m *SimpleMetrics) RecordError(err error) {
	m.errors = append(m.errors, err)
	log.Printf("ERROR: %v", err)
}

// Report prints a summary of metrics.
func (m *SimpleMetrics) Report() {
	log.Printf("=== Metrics Summary ===")
	log.Printf("Total events: %d", m.totalEvents)
	log.Printf("Total batches: %d", m.totalBatches)
	log.Printf("Errors: %d", len(m.errors))
	if len(m.errors) > 0 {
		log.Printf("First error: %v", m.errors[0])
	}
}
