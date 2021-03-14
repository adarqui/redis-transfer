package main

import (
	"log"
	"strconv"

	flag "github.com/spf13/pflag"
)

func init() {
	totalKeyCount = make(chan int, 1)
}

func main() {
	var replace bool
	var help bool
	flag.BoolVarP(&replace, "replace", "r", false, "whether replace existed keys")
	flag.BoolVarP(&help, "help", "h", false, "help")
	flag.Parse()

	args := flag.Args()
	if len(args) < 4 || help {
		usage()
	}

	from := args[0]
	to := args[1]
	keys := args[2]
	threads, err := strconv.Atoi(args[3])
	if err != nil {
		log.Fatal("Main: threads conversion error: ", err)
	}

	RunTransferArgs(from, to, keys, threads, replace)
}
