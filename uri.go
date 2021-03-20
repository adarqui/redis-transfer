package main

import (
	"fmt"
	goredis "gopkg.in/redis.v4"
	"log"
	"net"
	"strconv"
	"strings"
)

func parseURI(host string) (server *RedisServer, err error) {
	if strings.HasPrefix(host, "redis://") {
		server, err = parseRedisURI(host)
	} else {
		server, err = rHostSplit(host)
	}
	return
}

func rHostSplit(host string) (*RedisServer, error) {
	var tokens []string
	if strings.HasPrefix(host, "[") {
		// ipv6
		host = host[1:]
		tokens = strings.SplitN(host, "]:", 2)
		if len(tokens) != 2 {
			log.Fatal("rHostSplit: Needs <[host]:port> for IPv6 host")
		}
		host = tokens[0]
		tokens = strings.Split(tokens[1], ":")
	} else {
		// IPv4
		tokens = strings.Split(host, ":")
		if len(tokens) < 2 {
			log.Fatal("rHostSplit: Needs <host:port[:dbnum:[pass]]>")
		}
		host = tokens[0]
		tokens = tokens[1:]
	}
	port, err := strconv.Atoi(tokens[0])
	if err != nil {
		log.Fatal("rHostSplit: port conversion error: ", err)
	}

	serv := new(RedisServer)
	serv.host = host
	serv.port = port

	lenTokens := len(tokens)
	if lenTokens > 1 {
		db, err := strconv.Atoi(tokens[1])
		if err != nil {
			log.Fatal("rHostSplit: db conversion error: ", err)
		}
		serv.db = db
	}

	if lenTokens > 2 {
		serv.pass = tokens[2]
	}

	return serv, nil
}

func rHostCopy(r *RedisServer) (*RedisServer, error) {
	opts := &goredis.Options{
		Addr:     fmt.Sprintf(net.JoinHostPort(r.host, strconv.Itoa(r.port))),
		Password: r.pass,
		DB:       r.db,
	}
	c := goredis.NewClient(opts)
	rs := &RedisServer{
		client: c,
		host:   r.host,
		port:   r.port,
		db:     r.db,
		pass:   r.pass,
	}
	return rs, nil
}

func redisToString(s *RedisServer) string {
	return fmt.Sprintf("<redis://%s?db=%d>", net.JoinHostPort(s.host, strconv.Itoa(s.port)), s.db)
}
