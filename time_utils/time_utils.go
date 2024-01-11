package time_utils

import "time"

func IsToday(t time.Time) bool {
	tY, tM, tD := t.Date()
	nowY, nowM, nowD := time.Now().Date()
	return tD == nowD && tM == nowM && tY == nowY
}
