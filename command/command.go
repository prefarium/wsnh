package command

import (
	"errors"
	"wsnh/adapter"
)

const (
	CmdStart = "start"
	CmdStop  = "stop"
	CmdToday = "today"
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
	}

	return Command{}, errors.New("wrong command")
}
