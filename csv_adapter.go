package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type csvAdapter struct {
	csvPath string
}

func (a csvAdapter) openCSV() (*os.File, error) {
	f, openErr := os.Open(a.csvPath)

	if openErr == nil || !os.IsNotExist(openErr) {
		return f, openErr
	}

	if errDir := os.MkdirAll(filepath.Dir(a.csvPath), 0777); errDir != nil {
		return nil, errDir
	}

	return os.OpenFile(a.csvPath, os.O_RDONLY|os.O_CREATE, 0777)
}

func (a csvAdapter) readLast() (*entry, error) {
	f, openErr := a.openCSV()
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

func (a csvAdapter) readAll() ([]*entry, error) {
	f, openErr := os.Open(a.csvPath)
	if os.IsNotExist(openErr) {
		return make([]*entry, 0), nil
	} else if openErr != nil {
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

func (a csvAdapter) write(e *entry) error {
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

func (a csvAdapter) entryToCSV(e *entry) []string {
	return []string{string(e.command), e.timestamp.Format(time.RFC822Z)}
}

func (a csvAdapter) csvToEntry(line []string) (*entry, error) {
	cmd := command(line[0])
	timestamp, err := time.Parse(time.RFC822Z, line[1])
	if err != nil {
		return nil, fmt.Errorf("csv to entry conversion failure: %s", err)
	}

	return &entry{cmd, timestamp}, err
}
