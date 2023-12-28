package main

import (
	"errors"
	"fmt"
	"os"
	"time"
	"wsnh/adapters"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no command")
		return
	}

	cmd := command(os.Args[1])
	currentAdapter := adapters.CSVAdapter{CSVPath: "./db/entries.csv"}

	if cmd.isStart() {
		if err := trackTime(currentAdapter); err != nil {
			fmt.Println(err)
		}
		return
	}

	if cmd.isStop() {
		if err := stopTracking(currentAdapter); err != nil {
			fmt.Println(err)
		}
		return
	}

	if cmd == "today" {
		workedTime, err := timeWorkedToday(currentAdapter)

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

	if lastEntry != nil && command(lastEntry.Kind).isStart() {
		return errors.New("time is already ticking")
	}

	return a.Write(&adapters.Entry{Kind: CmdStart, Timestamp: time.Now()})
}

func stopTracking(a adapter) error {
	lastEntry, readErr := a.ReadLast()
	if readErr != nil {
		return readErr
	}

	if lastEntry == nil || command(lastEntry.Kind).isStop() {
		return errors.New("time is not being tracked")
	}

	return a.Write(&adapters.Entry{Kind: CmdStop, Timestamp: time.Now()})
}

func timeWorkedToday(a adapter) (time.Duration, error) {
	entries, readErr := a.ReadAll()
	if readErr != nil {
		return 0, readErr
	}

	tY, tM, tD := time.Now().Date()
	var lastStart time.Time
	var workedTime time.Duration

	for _, e := range entries {
		switch command(e.Kind) {
		case CmdStart:
			y, m, d := e.Timestamp.Date()
			if y == tY && m == tM && d == tD {
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
