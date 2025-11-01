package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdcb/cdc01/src/ingestors/telematics"
)

func main() {
	// CLI flags
	inputPath := flag.String("input", "", "Path to JSON Lines file or directory (required)")
	outputPath := flag.String("output", "", "Path to output file (optional, defaults to stdout)")
	batchSize := flag.Int("batch-size", 100, "Number of events per batch")
	verbose := flag.Bool("verbose", false, "Enable verbose output (prints all events)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Telematics Data Ingestor - Reads JSON Lines telematics events\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s --input data/raw/telematics --batch-size 50\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --input events.jsonl --output processed.jsonl\n", os.Args[0])
	}

	flag.Parse()

	if *inputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Validate input exists
	if _, err := os.Stat(*inputPath); err != nil {
		log.Fatalf("Input path does not exist: %s", *inputPath)
	}

	// Create output handler
	var outputHandler telematics.OutputHandler
	var err error

	if *outputPath != "" {
		outputHandler, err = telematics.NewFileHandler(*outputPath)
		if err != nil {
			log.Fatalf("Failed to create file handler: %v", err)
		}
		log.Printf("Output will be written to: %s", *outputPath)
	} else {
		outputHandler = telematics.NewStdoutHandler(*verbose)
		log.Printf("Output will be written to stdout")
	}

	// Create metrics handler
	metrics := telematics.NewSimpleMetrics()

	// Create ingestor
	ingestor, err := telematics.NewIngestor(telematics.IngestorConfig{
		InputPath:      *inputPath,
		BatchSize:      *batchSize,
		OutputHandler:  outputHandler,
		MetricsHandler: metrics,
	})
	if err != nil {
		log.Fatalf("Failed to create ingestor: %v", err)
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v, shutting down...", sig)
		cancel()
	}()

	// Run ingestor
	log.Printf("Starting ingestor (input=%s, batch_size=%d)", *inputPath, *batchSize)
	if err := ingestor.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("Ingestor error: %v", err)
	}

	// Print final metrics
	stats := ingestor.Stats()
	log.Printf("Ingestion complete:")
	log.Printf("  Events processed: %d", stats.EventsProcessed)
	log.Printf("  Batches processed: %d", stats.BatchesProcessed)
	log.Printf("  Errors: %d", stats.ErrorsEncountered)

	metrics.Report()
}
