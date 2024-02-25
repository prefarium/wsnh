package command

import (
	"time"
	"wsnh/time_utils"
)

func calcWeekTime(ds DataSource) (string, error) {
	entries, readErr := ds.ReadAll()
	if readErr != nil {
		return "", readErr
	}

	var (
		timeNow       = time.Now()
		weekBeginning = time_utils.BeginningOfWeek(timeNow)
		lastStart     time.Time
		workedTime    time.Duration
	)

	for _, e := range entries {
		switch e.Kind {
		case CmdStart:
			if e.Timestamp.Compare(weekBeginning) != -1 && e.Timestamp.Compare(timeNow) != 1 {
				lastStart = e.Timestamp
			}
		case CmdStop:
			if !lastStart.IsZero() {
				workedTime += e.Timestamp.Sub(lastStart)
				lastStart = time.Time{}
			}
		}
	}

	if !lastStart.IsZero() {
		workedTime += time.Now().Sub(lastStart)
	}

	return formatDuration(workedTime), nil
}
