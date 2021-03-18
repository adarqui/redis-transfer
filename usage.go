package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Println(
		`usage: redis-transfer <from_redis> <to_redis> <regex_or_input_file> <concurrent_threads> [--replace]

 * There are two redis connection formats to choose from:
   - host:port:[db[:password]]
   - redis://user:password@host:port?db=number

 * regex_or_input_file:
   - To transfer only those keys that match a regex pattern:
     some_key_prefix*
     *some_key_suffix
     prefix*something*suffix
     etc..

   - To transfer keys from an input file, simply list keys in a file, separated by newlines
     key1
     keyN

 * concurrent_threads:
   - This should be a number between 1 and (max_cpu's*10)
     You can play around with this to find the optimal setting. I generally use 50 on my 8 core box.

 * replace:
   - You can replace existed keys in the destination redis using --replace flag

 * examples:
   redis-transfer localhost:6379 remotehost:6379:1:password "migrate:*" 50
   redis-transfer redis://localhost:6379 redis://user:password@remotehost?db=1 "migrate:*" 50
   redis-transfer localhost:6379 remotehost:6379 "migrate:*" 50 --replace`)

	os.Exit(1)
}
