package command

import (
	"fmt"
	"strings"
	"time"
	"wsnh/time_utils"
)

func calcWeekTime(ds DataSource) (string, error) {
	entries, readErr := ds.ReadAll()
	if readErr != nil {
		return "", readErr
	}

	type workDay struct {
		weekday    time.Weekday
		workedTime time.Duration
	}

	var (
		timeNow       = time.Now()
		weekBeginning = time_utils.BeginningOfWeek(timeNow)
		windowStart   time.Time
		lastWorkedDay time.Time
		workedTime    time.Duration
		workedByDays  = make([]*workDay, 0, 7)
		output        strings.Builder
		workedTotal   time.Duration
	)

	for _, e := range entries {
		switch e.Kind {
		case CmdStart:
			if time_utils.IsCovered(e.Timestamp, weekBeginning, timeNow) {
				windowStart = e.Timestamp

				if !lastWorkedDay.IsZero() && !time_utils.IsSameDay(lastWorkedDay, windowStart) {
					workedByDays = append(workedByDays, &workDay{
						weekday:    lastWorkedDay.Weekday(),
						workedTime: workedTime,
					})
					workedTime = 0
				}

				lastWorkedDay = windowStart
			}
		case CmdStop:
			if !windowStart.IsZero() {
				workedInWindow := e.Timestamp.Sub(windowStart)
				workedTime += workedInWindow
				workedTotal += workedInWindow
				windowStart = time.Time{}
			}
		}
	}

	if !windowStart.IsZero() {
		workedTime += timeNow.Sub(windowStart)
		workedTotal += timeNow.Sub(windowStart)
	}

	if !lastWorkedDay.IsZero() && !time_utils.IsSameDay(lastWorkedDay, windowStart) {
		workedByDays = append(workedByDays, &workDay{
			weekday:    lastWorkedDay.Weekday(),
			workedTime: workedTime,
		})
	}

	for _, d := range workedByDays {
		output.WriteString(fmt.Sprintf("%s %s\n", formatDuration(d.workedTime), d.weekday))
	}

	output.WriteString(formatDuration(workedTotal))

	return output.String(), nil
}
