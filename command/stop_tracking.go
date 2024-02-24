package command

import (
	"errors"
	"time"
	"wsnh/adapter"
)

func stopTracking(ds DataSource) (string, error) {
	if lastEntry, err := ds.ReadLast(); err != nil {
		return "", err
	} else if lastEntry == nil || lastEntry.Kind == CmdStop {
		return "", errors.New("time is not being tracked")
	}

	if err := ds.Write(&adapter.Entry{Kind: CmdStop, Timestamp: time.Now()}); err != nil {
		return "", err
	}

	return "", nil
}
