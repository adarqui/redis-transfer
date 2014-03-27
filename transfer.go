package main

import (
	"github.com/cheggaaa/pb"
	"io/ioutil"
	"log"
	"menteslibres.net/gosexy/redis"
	"os"
	"strconv"
	"strings"
)

type Redis_Pipe struct {
	from    *Redis_Server
	to      *Redis_Server
	threads int
	keys    string
}

type Redis_Server struct {
	r    *redis.Client
	host string
	port int
	db   string
	user string
	pass string
}

type Op struct {
	str   string
	code  uint8
	repch chan bool
}

func usage() {
	log.Fatal("usage: ./transfer <from_redis_host:port[:dbNum[:pass]]> <to_redis_host:port[:dbNum[:pass]]> <key-regex or input-file-full-of-keys> <number-of-concurrent-threads>")
}

func rhost_split(host string) (*Redis_Server, error) {
	tokens := strings.Split(host, ":")
	if len(tokens) < 2 {
		log.Fatal("rhost_split: Needs <host:port[:dbnum:[pass]]>")
	}

	host = tokens[0]
	port, err := strconv.Atoi(tokens[1])
	if err != nil {
		log.Fatal("rhost_split: port conversion error: ", err)
	}

	serv := new(Redis_Server)
	serv.host = host
	serv.port = port

	len_tokens := len(tokens)
	if len_tokens > 2 {
		serv.db = tokens[2]
	}

	if len_tokens > 3 {
		serv.pass = tokens[3]
	}

	return serv, nil
}

func rhost_copy(r *Redis_Server) (*Redis_Server, error) {
	rnew := new(Redis_Server)
	rnew.r = redis.New()
	rnew.host = r.host
	rnew.port = r.port
	rnew.db = r.db
	rnew.pass = r.pass
	return rnew, nil
}

func New(from, to, keys string, threads int) *Redis_Pipe {
	pipe := new(Redis_Pipe)

	pipe.from, _ = rhost_split(from)
	pipe.to, _ = rhost_split(to)
	pipe.keys = keys

	pipe.threads = threads

	log.Printf("from=<%s:%i>, to=<%s:%i>\n", pipe.from.host, pipe.from.port, pipe.to.host, pipe.to.port)

	return pipe
}

func (pipe *Redis_Pipe) TransferThread(i int, ch chan Op) {
	log.Printf("transfer thread #%d launched\n", i)
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
	if serv.db != "" {
		dbnum, err2 := strconv.Atoi(serv.db)
		if err2 != nil {
			log.Fatal("ConnectOne: db number conversion error: ", err2)
		}
		_, err = serv.r.Select(int64(dbnum))
		if err != nil {
			log.Fatal("ConnectOne: select db failure: ", err)
		}
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

func (pipe *Redis_Pipe) KeysFile() []string {
	blob, err := ioutil.ReadFile(pipe.keys)
	if err != nil {
		log.Fatal("KeysFile: error reading keys file: ", err)
	}
	lines := strings.Split(string(blob), "\n")
	return lines
}

func (pipe *Redis_Pipe) KeysRedis() []string {
	keys, err := pipe.from.r.Keys(pipe.keys)
	if err != nil {
		log.Fatal("KeysRedis: error obtaining keys list from redis: ", err)
	}
	return keys
}

func (pipe *Redis_Pipe) Keys() []string {
	_, err := os.Stat(pipe.keys)

	var keys []string
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

	all_keys := pipes[0].Keys()

	count := len(all_keys)
	bar := pb.StartNew(count)
	bar.ShowPercent = true
	bar.ShowBar = true
	bar.ShowCounters = true
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true

	for _, v := range all_keys {
		op := Op{v, 0, nil}
		ch <- op
		bar.Increment()
	}

	for i := 0; i < pipe.threads; i++ {
		repch := make(chan bool, 1)
		op := Op{"", 1, repch}
		ch <- op
		_ = <-repch
	}

	bar.FinishPrint("Done.")

}
