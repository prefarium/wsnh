package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"wsnh/adapter"
)

func (a Adapter) ReadLast() (*adapter.Entry, error) {
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

func (a Adapter) openCSV() (*os.File, error) {
	f, openErr := os.Open(a.FilePath)

	if openErr == nil || !os.IsNotExist(openErr) {
		return f, openErr
	}

	if errDir := os.MkdirAll(filepath.Dir(a.FilePath), 0777); errDir != nil {
		return nil, errDir
	}

	return os.OpenFile(a.FilePath, os.O_RDONLY|os.O_CREATE, 0777)
}
