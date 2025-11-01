package telematics

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/mdcb/cdc01/pkg/models"
)

// IngestorConfig configures the telematics ingestor.
type IngestorConfig struct {
	// InputPath is the path to JSON Lines file(s)
	InputPath string
	// BatchSize controls how many events to read at once
	BatchSize int
	// OutputHandler processes batches of events
	OutputHandler OutputHandler
	// MetricsHandler receives metrics updates (optional)
	MetricsHandler MetricsHandler
}

// OutputHandler processes batches of telematics events.
// Implementations might write to Kafka, Redpanda, S3, or other sinks.
type OutputHandler interface {
	ProcessBatch(ctx context.Context, events []models.TelematicsEvent) error
	Close() error
}

// MetricsHandler receives metrics from the ingestor.
type MetricsHandler interface {
	RecordEventsProcessed(count int64)
	RecordBatchProcessed(durationMs int64)
	RecordError(err error)
}

// Ingestor reads telematics events and processes them through configured handlers.
type Ingestor struct {
	config         IngestorConfig
	reader         *Reader
	eventsRead     atomic.Int64
	batchesRead    atomic.Int64
	errorsEncountered atomic.Int64
}

// NewIngestor creates a new telematics ingestor.
func NewIngestor(cfg IngestorConfig) (*Ingestor, error) {
	if cfg.InputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}
	if cfg.OutputHandler == nil {
		return nil, fmt.Errorf("output handler is required")
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 100
	}

	reader, err := NewReader(ReaderConfig{
		InputPath: cfg.InputPath,
		BatchSize: cfg.BatchSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}

	return &Ingestor{
		config: cfg,
		reader: reader,
	}, nil
}

// Run starts the ingestor and processes events until completion or error.
func (i *Ingestor) Run(ctx context.Context) error {
	log.Printf("Starting telematics ingestor (input=%s, batch_size=%d)",
		i.config.InputPath, i.config.BatchSize)

	eventCh, errCh := i.reader.ReadEvents(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err, ok := <-errCh:
			if ok && err != nil {
				i.errorsEncountered.Add(1)
				if i.config.MetricsHandler != nil {
					i.config.MetricsHandler.RecordError(err)
				}
				return fmt.Errorf("reader error: %w", err)
			}

		case batch, ok := <-eventCh:
			if !ok {
				// Channel closed, all events processed
				return i.finalize()
			}

			if err := i.processBatch(ctx, batch); err != nil {
				i.errorsEncountered.Add(1)
				if i.config.MetricsHandler != nil {
					i.config.MetricsHandler.RecordError(err)
				}
				return fmt.Errorf("failed to process batch: %w", err)
			}
		}
	}
}

// processBatch processes a single batch of events.
func (i *Ingestor) processBatch(ctx context.Context, batch []models.TelematicsEvent) error {
	start := time.Now()

	if err := i.config.OutputHandler.ProcessBatch(ctx, batch); err != nil {
		return err
	}

	// Update metrics
	i.eventsRead.Add(int64(len(batch)))
	i.batchesRead.Add(1)

	if i.config.MetricsHandler != nil {
		i.config.MetricsHandler.RecordEventsProcessed(int64(len(batch)))
		i.config.MetricsHandler.RecordBatchProcessed(time.Since(start).Milliseconds())
	}

	return nil
}

// finalize completes the ingestion process and reports final stats.
func (i *Ingestor) finalize() error {
	if err := i.config.OutputHandler.Close(); err != nil {
		return fmt.Errorf("failed to close output handler: %w", err)
	}

	log.Printf("Ingestor completed: events=%d, batches=%d, errors=%d",
		i.eventsRead.Load(), i.batchesRead.Load(), i.errorsEncountered.Load())
	return nil
}

// Stats returns current ingestion statistics.
func (i *Ingestor) Stats() IngestorStats {
	return IngestorStats{
		EventsProcessed: i.eventsRead.Load(),
		BatchesProcessed: i.batchesRead.Load(),
		ErrorsEncountered: i.errorsEncountered.Load(),
	}
}

// IngestorStats holds ingestion statistics.
type IngestorStats struct {
	EventsProcessed   int64
	BatchesProcessed  int64
	ErrorsEncountered int64
}
