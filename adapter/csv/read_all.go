package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"wsnh/adapter"
)

func (a Adapter) ReadAll() ([]*adapter.Entry, error) {
	f, openErr := os.Open(a.FilePath)

	if os.IsNotExist(openErr) {
		return make([]*adapter.Entry, 0), nil
	}

	if openErr != nil {
		return nil, fmt.Errorf("failed to open csv: %s", openErr)
	}

	defer f.Close()

	lines, readErr := csv.NewReader(f).ReadAll()
	if readErr != nil {
		return nil, fmt.Errorf("failed to read csv: %s", readErr)
	}

	entries := make([]*adapter.Entry, len(lines))
	for i, line := range lines {
		e, parseErr := a.csvToEntry(line)
		if parseErr != nil {
			return entries, parseErr
		}

		entries[i] = e
	}

	return entries, nil
}
