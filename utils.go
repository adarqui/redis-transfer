package main

import (
	"log"
	"net"
	"net/url"
	"strconv"
)

func parseRedisURI(s string) (server *Redis_Server, err error) {
	// "redis://x:password@host.name:123?db=0"
	// "redis://x:password@host.name:123"

	// Defaults
	host := "localhost"
	password := ""
	port := 6379
	db := 0

	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	if u.Scheme != "redis" {
		log.Fatal("Scheme must be redis")
	}
	q := u.Query()
	dbS := q.Get("db")
	if u.User != nil {
		var ok bool
		password, ok = u.User.Password()
		if !ok {
			password = ""
		}
	}

	var p string
	host, p, _ = net.SplitHostPort(u.Host)

	if p != "" {
		port, err = strconv.Atoi(p)
		if err != nil {
			log.Fatalf("Unable to convert port to integer for %s", err)
		}
	}

	if dbS != "" {
		db, err = strconv.Atoi(dbS)
		if err != nil {
			log.Fatalf("Unable to convert db to integer for %s", dbS)
		}
	}

	return &Redis_Server{host: host, port: port, db: db, pass: password}, nil
}
