package main

import (
	goredis "gopkg.in/redis.v4"
	"sync"
)

type Redis_Pipe struct {
	from    *RedisServer
	to      *RedisServer
	threads int
	keys    string
}

type RedisServer struct {
	client *goredis.Client
	host   string
	port   int
	db     int
	pass   string
}

const (
	OP_NOP = 0
	OP_DIE = iota
)

type Op struct {
	str   string
	code  uint8
	repch chan bool
}

type redisKey string

var totalKeyCount chan int

var wg sync.WaitGroup
