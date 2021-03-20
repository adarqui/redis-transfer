package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertFieldsEqual(t *testing.T, e *RedisServer, a *RedisServer) {
	assert.Equal(t, e.host, a.host, "host should be set")
	assert.Equal(t, e.port, a.port, "port should be set")
	assert.Equal(t, e.db, a.db, "db should be set")
	assert.Equal(t, e.pass, a.pass, "password should be set")
}

func TestParseURI(t *testing.T) {
	s := "redis://x:password@host.com:123"
	a, _ := parseURI(s)
	e := &RedisServer{
		host: "host.com",
		port: 123,
		db:   0,
		pass: "password",
	}
	assertFieldsEqual(t, e, a)

	s2 := "host.com:123:0:password"
	a2, _ := parseURI(s2)
	assertFieldsEqual(t, e, a2)

	s3 := "redis://localhost:6370"
	a3, _ := parseURI(s3)
	e3 := &RedisServer{nil, "localhost", 6370, 0, ""}
	assertFieldsEqual(t, e3, a3)
}

func TestRedisURIParsing(t *testing.T) {
	s := "redis://x:password@host.com:123"
	a, _ := parseRedisURI(s)
	e := &RedisServer{
		host: "host.com",
		port: 123,
		db:   0,
		pass: "password",
	}
	assertFieldsEqual(t, e, a)
}

func TestRedisURIParsingWithDB(t *testing.T) {
	s := "redis://x:password@host.com:123?db=0"
	actual, _ := parseRedisURI(s)
	expected := &RedisServer{
		host: "host.com",
		port: 123,
		db:   0,
		pass: "password",
	}
	assertFieldsEqual(t, expected, actual)
}

func TestRHost_Split(t *testing.T) {
	s := "localhost:6370:0:password"
	actual, _ := rHostSplit(s)
	expected := &RedisServer{nil, "localhost", 6370, 0, "password"}

	assertFieldsEqual(t, expected, actual)
}
