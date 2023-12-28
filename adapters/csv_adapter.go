package adapters

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	TimeFormat = time.RFC822Z
)

type CSVAdapter struct {
	CSVPath string
}

func (a CSVAdapter) ReadLast() (*Entry, error) {
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
		}

		if errors.Is(err, io.EOF) {
			break
		}

		return nil, fmt.Errorf("failed to read csv: %s", err)
	}

	if lastLine == nil {
		return nil, nil
	}

	return a.csvToEntry(lastLine)
}

func (a CSVAdapter) ReadAll() ([]*Entry, error) {
	f, openErr := os.Open(a.CSVPath)

	if os.IsNotExist(openErr) {
		return make([]*Entry, 0), nil
	}

	if openErr != nil {
		return nil, fmt.Errorf("failed to open csv: %s", openErr)
	}

	defer f.Close()

	lines, readErr := csv.NewReader(f).ReadAll()
	if readErr != nil {
		return nil, fmt.Errorf("failed to read csv: %s", readErr)
	}

	entries := make([]*Entry, len(lines))
	for i, line := range lines {
		e, parseErr := a.csvToEntry(line)
		if parseErr != nil {
			return entries, parseErr
		}

		entries[i] = e
	}

	return entries, nil
}

func (a CSVAdapter) Write(e *Entry) error {
	f, openErr := os.OpenFile(a.CSVPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if openErr != nil {
		return fmt.Errorf("failed to open csv: %s", openErr)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if writeErr := w.Write(a.entryToCSV(e)); writeErr != nil {
		return fmt.Errorf("failed to write csv: %s", writeErr)
	}

	w.Flush()
	return nil
}

func (a CSVAdapter) openCSV() (*os.File, error) {
	f, openErr := os.Open(a.CSVPath)

	if openErr == nil || !os.IsNotExist(openErr) {
		return f, openErr
	}

	if errDir := os.MkdirAll(filepath.Dir(a.CSVPath), 0777); errDir != nil {
		return nil, errDir
	}

	return os.OpenFile(a.CSVPath, os.O_RDONLY|os.O_CREATE, 0777)
}

func (a CSVAdapter) entryToCSV(e *Entry) []string {
	kind := e.Kind
	timestamp := e.Timestamp.Format(TimeFormat)

	return []string{kind, timestamp}
}

func (a CSVAdapter) csvToEntry(line []string) (*Entry, error) {
	kind := line[0]
	timestamp, err := time.Parse(TimeFormat, line[1])

	if err != nil {
		return nil, fmt.Errorf("csv to entry conversion failure: %s", err)
	}

	return &Entry{kind, timestamp}, err
}
