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
	ad := adapters.CSVAdapter{CSVPath: "./db/entries.csv"}

	if cmd.isStart() {
		err := startTracking(ad)
		if err != nil {
			fmt.Println(err)
		}
	} else if cmd.isStop() {
		err := stopTracking(ad)
		if err != nil {
			fmt.Println(err)
		}
	} else if cmd == "today" {
		workedTime, err := calcToday(ad)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(workedTime)
		}
	} else {
		fmt.Println("wrong command")
	}
}

func startTracking(a adapter) error {
	lastEntry, readErr := a.ReadLast()

	if readErr != nil {
		return readErr
	} else if lastEntry != nil && command(lastEntry.Kind).isStart() {
		return errors.New("time is already ticking")
	} else {
		writeErr := a.Write(&adapters.Entry{Kind: CmdStart, Timestamp: time.Now()})
		if writeErr != nil {
			return writeErr
		} else {
			return nil
		}
	}
}

func stopTracking(a adapter) error {
	lastEntry, readErr := a.ReadLast()

	if readErr != nil {
		return readErr
	} else if lastEntry == nil || command(lastEntry.Kind).isStop() {
		return errors.New("time is not being tracked")
	} else {
		writeErr := a.Write(&adapters.Entry{Kind: CmdStop, Timestamp: time.Now()})
		if writeErr != nil {
			return writeErr
		} else {
			return nil
		}
	}
}

func calcToday(a adapter) (time.Duration, error) {
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
		lastStart = time.Time{}
	}

	return workedTime, nil
}
