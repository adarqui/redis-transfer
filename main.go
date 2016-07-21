package main

import (
	"log"
	"os"
	"strconv"
)

func init() {
	totalKeyCount = make(chan int, 1)
}

func main() {

	if len(os.Args) < 5 {
		usage()
	}

	from := os.Args[1]
	to := os.Args[2]
	keys := os.Args[3]
	threads, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Fatal("Main: threads conversion error: ", err)
	}

	RunTransferArgs(from, to, keys, threads)
}
