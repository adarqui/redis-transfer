package main

import (
	goredis "gopkg.in/redis.v4"
	"menteslibres.net/gosexy/redis"
	"sync"
)

type Redis_Pipe struct {
	from    *Redis_Server
	to      *Redis_Server
	threads int
	keys    string
}

type Redis_Server struct {
	r      *redis.Client
	client *goredis.Client
	host   string
	port   int
	db     int
	pass   string
}

type Op struct {
	str   string
	code  uint8
	repch chan bool
}

type redisKey string

var totalKeyCount chan int

var wg sync.WaitGroup
