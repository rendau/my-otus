package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	argFrom   string
	argTo     string
	argOffset int64
	argLimit  int64
)

func init() {
	flag.StringVar(&argFrom, "from", "", "file to read from (required)")
	flag.StringVar(&argTo, "to", "", "file to write to (required)")
	flag.Int64Var(&argOffset, "offset", 0, "bytes to omit in source file, before copy (default 0)")
	flag.Int64Var(&argLimit, "limit", 0, "bytes count to copy (default 0, that mean 'no limit')")
}

func main() {
	flag.Parse()

	// Validation arguments
	if argFrom == "" || argTo == "" {
		flag.Usage()
		os.Exit(1)
	}
	if argOffset < 0 {
		fmt.Println("offset must be a positive number")
		flag.Usage()
		os.Exit(1)
	}
	if argLimit < 0 {
		fmt.Println("limit must be a positive number")
		flag.Usage()
		os.Exit(1)
	}

	err := fileCopy(argFrom, argTo, argOffset, argLimit)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
