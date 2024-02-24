package command

import (
	"time"
	"wsnh/time_utils"
)

func calcTodayTime(ds DataSource) (string, error) {
	entries, readErr := ds.ReadAll()
	if readErr != nil {
		return "", readErr
	}

	var lastStart time.Time
	var workedTime time.Duration

	for _, e := range entries {
		switch e.Kind {
		case CmdStart:
			if time_utils.IsToday(e.Timestamp) {
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

	return workedTime.String(), nil
}
