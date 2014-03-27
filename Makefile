all:
	go build

deps:
	go get github.com/cheggaaa/pb
	go get menteslibres.net/gosexy/redis

clean:
	rm -f redis-transfer
