package main

import "time"

type entry struct {
	command   command
	timestamp time.Time
}
