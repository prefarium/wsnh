package csv

import (
	"fmt"
	"time"
	"wsnh/adapter"
)

type Adapter struct {
	FilePath   string
	TimeFormat string
}

func NewAdapter(filePath string) *Adapter {
	return &Adapter{
		FilePath:   filePath,
		TimeFormat: time.RFC822Z,
	}
}

func (a Adapter) csvToEntry(line []string) (*adapter.Entry, error) {
	kind := line[0]
	timestamp, err := time.Parse(a.TimeFormat, line[1])

	if err != nil {
		return nil, fmt.Errorf("csv to entry conversion failure: %s", err)
	}

	return &adapter.Entry{Kind: kind, Timestamp: timestamp}, err
}
