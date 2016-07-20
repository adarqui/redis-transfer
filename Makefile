build: deps
	go build

deps:
	go get github.com/cheggaaa/pb
	go get gopkg.in/redis.v4
	go get menteslibres.net/gosexy/redis

test:
	go test

clean:
	rm -f redis-transfer
