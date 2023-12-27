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

type CSVAdapter struct {
	CSVPath string
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

func (a CSVAdapter) ReadAll() ([]*Entry, error) {
	f, openErr := os.Open(a.CSVPath)
	if os.IsNotExist(openErr) {
		return make([]*Entry, 0), nil
	} else if openErr != nil {
		return nil, fmt.Errorf("failed to open csv: %s", openErr)
	}
	defer f.Close()

	r := csv.NewReader(f)
	lines, readErr := r.ReadAll()
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
	writeErr := w.Write(a.entryToCSV(e))
	if writeErr != nil {
		return fmt.Errorf("failed to write csv: %s", writeErr)
	} else {
		w.Flush()
		return nil
	}
}

func (a CSVAdapter) entryToCSV(e *Entry) []string {
	return []string{e.Kind, e.Timestamp.Format(time.RFC822Z)}
}

func (a CSVAdapter) csvToEntry(line []string) (*Entry, error) {
	kind := line[0]
	timestamp, err := time.Parse(time.RFC822Z, line[1])
	if err != nil {
		return nil, fmt.Errorf("csv to entry conversion failure: %s", err)
	}

	return &Entry{kind, timestamp}, err
}
