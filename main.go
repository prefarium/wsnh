package main

import (
	"errors"
	"fmt"
	"os"
	"time"
	"wsnh/adapters"
	"wsnh/time_utils"
)

const (
	CmdStart = "start"
	CmdStop  = "stop"
	CmdToday = "today"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no command")
		return
	}

	cmd := os.Args[1]
	a := adapters.NewCSVAdapter("./db/entries.csv")

	if cmd == CmdStart {
		if err := trackTime(a); err != nil {
			fmt.Println(err)
		}
		return
	}

	if cmd == CmdStop {
		if err := stopTracking(a); err != nil {
			fmt.Println(err)
		}
		return
	}

	if cmd == CmdToday {
		workedTime, err := timeWorkedToday(a)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(workedTime)
		}

		return
	}

	fmt.Println("wrong command")
}

func trackTime(a adapter) error {
	lastEntry, readErr := a.ReadLast()
	if readErr != nil {
		return readErr
	}

	if lastEntry != nil && lastEntry.Kind == CmdStart {
		return errors.New("time is already ticking")
	}

	return a.Write(&adapters.Entry{Kind: CmdStart, Timestamp: time.Now()})
}

func stopTracking(a adapter) error {
	lastEntry, readErr := a.ReadLast()
	if readErr != nil {
		return readErr
	}

	if lastEntry == nil || lastEntry.Kind == CmdStop {
		return errors.New("time is not being tracked")
	}

	return a.Write(&adapters.Entry{Kind: CmdStop, Timestamp: time.Now()})
}

func timeWorkedToday(a adapter) (time.Duration, error) {
	entries, readErr := a.ReadAll()
	if readErr != nil {
		return 0, readErr
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

	return workedTime, nil
}
