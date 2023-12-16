package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no command")
		return
	}

	cmd := command(os.Args[1])
	ad := CSVAdapter{"./db/entries.csv"}

	if cmd.isStart() {
		lastEntry, readErr := ad.readLast()
		if readErr != nil {
			fmt.Println(readErr)
		} else if lastEntry != nil && lastEntry.command.isStart() {
			fmt.Println("time is already ticking")
		} else {
			writeErr := ad.write(newEntry(cmd))
			if writeErr != nil {
				fmt.Println(writeErr)
			}
		}
	} else if cmd.isEnd() {
		lastEntry, readErr := ad.readLast()
		if readErr != nil {
			fmt.Println(readErr)
		} else if lastEntry == nil || lastEntry.command.isEnd() {
			fmt.Println("time is not being tracked")
		} else {
			writeErr := ad.write(newEntry(cmd))
			if writeErr != nil {
				fmt.Println(writeErr)
			}
		}
	} else {
		fmt.Println("wrong command")
	}
}

func newEntry(cmd command) *entry {
	return &entry{cmd, time.Now()}
}
