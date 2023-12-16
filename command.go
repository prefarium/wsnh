package main

type command string

func (c command) isStart() bool {
	return c == "start" || c == "go" || c == "begin"
}

func (c command) isEnd() bool {
	return c == "end" || c == "stop" || c == "finish"
}

func (c command) String() string {
	if c.isStart() {
		return "start"
	} else if c.isEnd() {
		return "end"
	} else {
		return string(c)
	}
}
