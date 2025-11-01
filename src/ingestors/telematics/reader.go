package telematics

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mdcb/cdc01/pkg/models"
)

// Reader reads telematics events from JSON Lines files.
type Reader struct {
	config ReaderConfig
}

// ReaderConfig configures the JSON Lines reader.
type ReaderConfig struct {
	// InputPath can be a file path or directory containing .jsonl files
	InputPath string
	// BatchSize controls how many events to read before yielding
	BatchSize int
}

// NewReader creates a new telematics event reader.
func NewReader(cfg ReaderConfig) (*Reader, error) {
	if cfg.InputPath == "" {
		return nil, fmt.Errorf("input path cannot be empty")
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 100 // default batch size
	}
	return &Reader{config: cfg}, nil
}

// ReadEvents reads telematics events from the configured input path.
// It yields batches of events via a channel and closes it when done.
func (r *Reader) ReadEvents(ctx context.Context) (<-chan []models.TelematicsEvent, <-chan error) {
	eventCh := make(chan []models.TelematicsEvent, 10)
	errCh := make(chan error, 1)

	go func() {
		defer close(eventCh)
		defer close(errCh)

		files, err := r.findJSONLFiles()
		if err != nil {
			errCh <- fmt.Errorf("failed to find input files: %w", err)
			return
		}

		for _, file := range files {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
			}

			if err := r.readFile(ctx, file, eventCh); err != nil {
				errCh <- fmt.Errorf("failed to read %s: %w", file, err)
				return
			}
		}
	}()

	return eventCh, errCh
}

// findJSONLFiles discovers all .jsonl files in the configured input path.
func (r *Reader) findJSONLFiles() ([]string, error) {
	info, err := os.Stat(r.config.InputPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		// Single file
		return []string{r.config.InputPath}, nil
	}

	// Directory - find all .jsonl files
	var files []string
	err = filepath.Walk(r.config.InputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".jsonl" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// readFile reads a single JSON Lines file and sends events in batches.
func (r *Reader) readFile(ctx context.Context, path string, eventCh chan<- []models.TelematicsEvent) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// Increase buffer size for large JSON lines
	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	batch := make([]models.TelematicsEvent, 0, r.config.BatchSize)
	lineNum := 0

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue // skip empty lines
		}

		var event models.TelematicsEvent
		if err := json.Unmarshal(line, &event); err != nil {
			return fmt.Errorf("line %d: failed to parse JSON: %w", lineNum, err)
		}

		batch = append(batch, event)

		// Send batch when full
		if len(batch) >= r.config.BatchSize {
			eventCh <- batch
			batch = make([]models.TelematicsEvent, 0, r.config.BatchSize)
		}
	}

	// Send remaining events
	if len(batch) > 0 {
		eventCh <- batch
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}
