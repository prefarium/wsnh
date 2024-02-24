package command

import (
	"errors"
	"time"
	"wsnh/adapter"
)

func startTracking(ds DataSource) (string, error) {
	if lastEntry, err := ds.ReadLast(); err != nil {
		return "", err
	} else if lastEntry != nil && lastEntry.Kind == CmdStart {
		return "", errors.New("time is already ticking")
	}

	if err := ds.Write(&adapter.Entry{Kind: CmdStart, Timestamp: time.Now()}); err != nil {
		return "", err
	}

	return "", nil
}
