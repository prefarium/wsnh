package time_utils

import "time"

func IsToday(t time.Time) bool {
	tY, tM, tD := t.Date()
	nowY, nowM, nowD := time.Now().Date()
	return tD == nowD && tM == nowM && tY == nowY
}

func BeginningOfWeek(t time.Time) time.Time {
	var d int64

	if wDay := int(t.Weekday()); wDay == 0 {
		d = 6
	} else {
		d = int64(wDay) - 1
	}

	return truncateToDay(t.Add(-time.Duration(d * 24 * int64(time.Hour))))
}

func IsSameDay(t1, t2 time.Time) bool {
	return truncateToDay(t1).Compare(truncateToDay(t2)) == 0
}

func IsCovered(t, start, end time.Time) bool {
	return t.Compare(start) != -1 && t.Compare(end) != 1
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
