package main

import (
	"fmt"
	goredis "gopkg.in/redis.v4"
	"log"
	"strconv"
	"strings"
)

func parseURI(host string) (server *Redis_Server, err error) {
	if strings.HasPrefix(host, "redis://") {
		server, err = parseRedisURI(host)
	} else {
		server, err = rhost_split(host)
	}
	return
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
		db, err := strconv.Atoi(tokens[2])
		if err != nil {
			log.Fatal("rhost_split: db conversion error: ", err)
		}
		serv.db = db
	}

	if len_tokens > 3 {
		serv.pass = tokens[3]
	}

	return serv, nil
}

func rhost_copy(r *Redis_Server) (*Redis_Server, error) {
	opts := &goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.host, r.port),
		Password: r.pass,
		DB:       r.db,
	}
	c := goredis.NewClient(opts)
	rs := &Redis_Server{
		client: c,
		host:   r.host,
		port:   r.port,
		db:     r.db,
		pass:   r.pass,
	}
	return rs, nil
}

func redisToString(s *Redis_Server) string {
	return fmt.Sprintf("<redis://%s:%d?db=%d>", s.host, s.port, s.db)
}
