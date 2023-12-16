package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type CSVAdapter struct {
	csvPath string
}

func (a CSVAdapter) readLast() (*entry, error) {
	f, openErr := os.Open(a.csvPath)
	if openErr != nil {
		return nil, fmt.Errorf("failed to open csv: %s", openErr)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.ReuseRecord = true
	var lastLine []string

	for {
		line, err := r.Read()
		if err == nil {
			lastLine = line
			continue
		} else if errors.Is(err, io.EOF) {
			break
		} else {
			return nil, fmt.Errorf("failed to read csv: %s", err)
		}
	}

	if lastLine == nil {
		return nil, nil
	}

	e, parseErr := a.csvToEntry(lastLine)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to read csv: %s", parseErr)
	}

	return e, nil
}

func (a CSVAdapter) readAll() ([]*entry, error) {
	f, openErr := os.Open(a.csvPath)
	if openErr != nil {
		return nil, fmt.Errorf("failed to open csv: %s", openErr)
	}
	defer f.Close()

	r := csv.NewReader(f)
	lines, readErr := r.ReadAll()
	if readErr != nil {
		return nil, fmt.Errorf("failed to read csv: %s", readErr)
	}

	entries := make([]*entry, len(lines))
	for i, line := range lines {
		timestamp, err := time.Parse(time.RFC822Z, line[1])
		if err != nil {
			return entries, fmt.Errorf("csv to entry conversion failure: %s", err)
		}

		entries[i] = &entry{command(line[0]), timestamp}
	}

	return entries, nil
}

func (a CSVAdapter) write(e *entry) error {
	f, openErr := os.OpenFile(a.csvPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if openErr != nil {
		return fmt.Errorf("failed to open csv: %s", openErr)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	writeErr := w.Write(a.entryToCSV(e))
	if writeErr != nil {
		return fmt.Errorf("failed to write csv: %s", writeErr)
	} else {
		w.Flush()
		return nil
	}
}

func (a CSVAdapter) entryToCSV(e *entry) []string {
	return []string{string(e.command), e.timestamp.Format(time.RFC822Z)}
}

func (a CSVAdapter) csvToEntry(line []string) (*entry, error) {
	cmd := command(line[0])
	timestamp, err := time.Parse(time.RFC822Z, line[1])
	if err != nil {
		return nil, fmt.Errorf("csv to entry conversion failure: %s", err)
	}

	return &entry{cmd, timestamp}, err
}
