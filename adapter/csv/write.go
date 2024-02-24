package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"wsnh/adapter"
)

func (a Adapter) Write(e *adapter.Entry) error {
	f, openErr := os.OpenFile(a.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
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

func (a Adapter) entryToCSV(e *adapter.Entry) []string {
	kind := e.Kind
	timestamp := e.Timestamp.Format(a.TimeFormat)

	return []string{kind, timestamp}
}
