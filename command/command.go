package command

import (
	"errors"
	"fmt"
	"time"
	"wsnh/adapter"
)

const (
	CmdStart = "start"
	CmdStop  = "stop"
	CmdToday = "today"
	CmdWeek  = "week"
)

type DataSource interface {
	ReadLast() (*adapter.Entry, error)
	ReadAll() ([]*adapter.Entry, error)
	Write(*adapter.Entry) error
}

type Command struct {
	repo DataSource
	Run  func(DataSource) (string, error)
}

func NewCommand(cmd string) (Command, error) {
	if cmd == CmdStart {
		return Command{Run: startTracking}, nil
	} else if cmd == CmdStop {
		return Command{Run: stopTracking}, nil
	} else if cmd == CmdToday {
		return Command{Run: calcTodayTime}, nil
	} else if cmd == CmdWeek {
		return Command{Run: calcWeekTime}, nil
	}

	return Command{}, errors.New("wrong command")
}

func formatDuration(d time.Duration) string {
	h := d / time.Hour
	m := (d - h*time.Hour) / time.Minute

	return fmt.Sprintf("%02dh%02dm", h, m)
}
