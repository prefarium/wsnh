package main

type command string

const (
	CmdStart = "start"
	CmdStop  = "stop"
)

func (c command) isStart() bool {
	return c == CmdStart || c == "go" || c == "begin"
}

func (c command) isStop() bool {
	return c == CmdStop || c == "finish" || c == "end"
}
