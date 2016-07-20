package main

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func New(from, to, keys string, threads int) *Redis_Pipe {
	pipe := new(Redis_Pipe)

	pipe.from, _ = parseURI(from)
	pipe.to, _ = parseURI(to)
	pipe.keys = keys

	pipe.threads = threads

	log.Printf("Initiating transfer from %s to %s\n", redisToString(pipe.from), redisToString(pipe.to))

	return pipe
}

func (pipe *Redis_Pipe) TransferThread(i int, ch chan Op) {
	for m := range ch {
		if m.code == 1 {
			// force children to exit, just reply true & vaporize this go routine
			m.repch <- true
			return
		}
		data, err := pipe.from.r.Dump(m.str)
		if err != nil {
			log.Printf("FAIL:DUMP:%s, %v\n", m.str, err)
		}
		if len(data) == 0 {
			continue
		}
		_, err = pipe.to.r.Restore(m.str, 0, data)
		if err != nil {
			log.Printf("FAIL:RESTORE:%s, %v\n", m.str, err)
		}
	}
}

func (serv *Redis_Server) ConnectOne() error {
	err := serv.r.ConnectNonBlock(serv.host, uint(serv.port))
	if err != nil {
		log.Fatal("ConnectOne: Connecting to host/port: ", err)
	}
	if serv.pass != "" {
		_, err = serv.r.Auth(serv.pass)
		if err != nil {
			log.Fatal("ConnectOne: pass incorrect: ", err)
		}
	}
	_, err = serv.r.Select(int64(serv.db))
	if err != nil {
		log.Fatal("ConnectOne: select db failure: ", err)
	}
	return nil
}

func (pipe *Redis_Pipe) Connect() error {
	err := pipe.from.ConnectOne()
	if err != nil {
		log.Fatal("Connect: Connecting to \"from\" host/port: ", err)
	}
	err = pipe.to.ConnectOne()
	if err != nil {
		log.Fatal("Connect: Connecting to \"to\" host/port: ", err)
	}
	return nil
}

func (pipe *Redis_Pipe) Init() ([]Redis_Pipe, chan Op) {

	pipes := make([]Redis_Pipe, pipe.threads)

	for i := 0; i < pipe.threads; i++ {
		pipes[i].keys = pipe.keys
		pipes[i].from, _ = rhost_copy(pipe.from)
		pipes[i].to, _ = rhost_copy(pipe.to)

		/* connect to both redii */
		pipes[i].Connect()
	}

	ch := make(chan Op, pipe.threads)

	for i := 0; i < pipe.threads; i++ {
		_i := i
		go pipes[_i].TransferThread(_i, ch)
	}

	return pipes, ch
}

func (pipe *Redis_Pipe) KeysFile() chan redisKey {
	blob, err := ioutil.ReadFile(pipe.keys)
	if err != nil {
		log.Fatal("KeysFile: error reading keys file: ", err)
	}
	keyChan := make(chan redisKey)
	lines := filter(strings.Split(string(blob), "\n"), func(s string) bool { return len(s) > 0 })
	totalKeyCount <- len(lines)
	go func() {
		for _, line := range lines {
			keyChan <- redisKey(line)
		}
	}()
	return keyChan
}

func init() {
	totalKeyCount = make(chan int, 1)
}

func (pipe *Redis_Pipe) KeysRedis() chan redisKey {
	keyChan := make(chan redisKey, 1000)
	info := pipe.from.client.Info("keyspace")
	// Sample: db0:keys=1201,expires=0,avg_ttl=0
	keyRegex := fmt.Sprintf("db%d:keys=(\\d+)", pipe.from.db)
	re := regexp.MustCompile(keyRegex)
	m := re.FindStringSubmatch(info.Val())
	if len(m) > 1 {
		if ks, err := strconv.Atoi(m[1]); err == nil {
			totalKeyCount <- ks
		}
	}
	split := make(chan []string)
	splitter := func() {
		wg.Add(1)
		defer wg.Done()
		defer close(keyChan)
		for ks := range split {
			for _, k := range ks {
				keyChan <- redisKey(k)
			}
		}
	}

	go splitter()

	go func(c chan redisKey) {
		wg.Add(1)
		defer wg.Done()
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			// REDIS SCAN
			// http://redis.io/commands/scan
			// Preferable because it doesn't lock complete database on larger keysets for 250ms+.
			keys, cursor, err = pipe.from.client.Scan(cursor, pipe.keys, 1000).Result()
			if err != nil {
				log.Fatal("SCAN: error obtaining key scan from redis: ", err)
			}
			split <- keys

			n += len(keys)
			if cursor == 0 {
				close(split)
				break
			}
		}
	}(keyChan)

	return keyChan
}

func (pipe *Redis_Pipe) Keys() chan redisKey {

	_, err := os.Stat(pipe.keys)

	var keys chan redisKey
	if err == nil {
		keys = pipe.KeysFile()
	} else {
		keys = pipe.KeysRedis()
	}

	return keys
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
		op := Op{string(v), 0, nil}
		ch <- op
		bar.Increment()
	}

	for i := 0; i < pipe.threads; i++ {
		repch := make(chan bool, 1)
		op := Op{"", 1, repch}
		ch <- op
		_ = <-repch
	}

	wg.Wait()

	bar.Finish()
	log.Println("Done.")
}
