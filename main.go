package main

import (
	"github.com/cheggaaa/pb"
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

	if threads <= 0 {
		log.Fatal("Main: threads must be > 0")
	}

	pipe := New(from, to, keys, threads)
	pipes, ch := pipe.Init()

	// Provides us with a channel that returns keys from redis or from a file
	keyChan := pipes[0].Keys()

	count := len(keyChan)
	bar := pb.StartNew(count)
	bar.ShowPercent = true
	bar.ShowBar = true
	bar.ShowCounters = true
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true

	wg.Add(1)
	go func() {
		defer wg.Done()
		t := <-totalKeyCount
		bar.Total = int64(t)
	}()

	for v := range keyChan {
		op := Op{string(v), OP_NOP, nil}
		ch <- op
		bar.Increment()
	}

	for i := 0; i < pipe.threads; i++ {
		repch := make(chan bool, 1)
		op := Op{"", OP_DIE, repch}
		ch <- op
		_ = <-repch
	}

	wg.Wait()

	bar.Finish()
	log.Println("Done.")
}
