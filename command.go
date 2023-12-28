package main

const (
	CmdStart = "start"
	CmdStop  = "stop"
)

type command string

func (c command) isStart() bool {
	return c == CmdStart || c == "go" || c == "begin"
}

func (c command) isStop() bool {
	return c == CmdStop || c == "finish" || c == "end"
}
