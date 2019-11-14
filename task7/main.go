package main

import (
	"fmt"
	"github.com/rendau/my-otus/task7/internal"
	"log"
	"os"
)

func main() {
	ok, d, cmd := parseArgs()
	if !ok {
		os.Exit(1)
	}

	envs, err := internal.ReadDir(d)
	if err != nil {
		log.Fatalln("Fail to run ReadDir")
	}

	exitCode := internal.RunCmd(cmd, envs, nil, nil)
	if exitCode == 111 {
		log.Fatalln("Fail to run RunCmd")
	}

	os.Exit(exitCode)
}

func parseArgs() (bool, string, []string) {
	args := os.Args

	if len(args) < 3 {
		fmt.Println(`Usage:
  task7 /path/to/env-dir command [arg1, arg2, ...]`)
		return false, "", nil
	}

	return true, args[1], args[2:]
}
