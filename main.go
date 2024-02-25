package main

import (
	"fmt"
	"os"
	"wsnh/adapter/csv"
	"wsnh/command"
	"wsnh/input"
)

func main() {
	cmd, parseErr := input.Parse(os.Args)
	if parseErr != nil {
		fmt.Println(parseErr)
	}

	runner, cmdErr := command.NewCommand(cmd)
	if cmdErr != nil {
		fmt.Println(cmdErr)
		return
	}

	if output, err := runner.Run(csv.NewAdapter("./db/database.csv")); err != nil {
		fmt.Println(err)
	} else if output != "" {
		fmt.Println(output)
	}
}
