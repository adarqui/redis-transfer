What it does?
"Parallelizes" transferring of redis keys from a source to a target host. You can specify what keys to transfer via regex ("*", "bleh:*") or file ("keys.txt","/tmp/keys.txt"). It will then spawn N go routines which pull (dump) from the source redis & immediately push (restore) to the target redis.


Usage:
$ ./redis-transfer
2014/03/26 23:19:03 usage: ./transfer <from_redis_host:port[:dbNum[:pass]]> <to_redis_host:port[:dbNum[:pass]]> <key-regex or input-file-full-of-keys> <number-of-threads>


Example usage:
./redis-transfer localhost:6379 localhost:6370:0:password "*" 40


./redis-transfer localhost:6379:2 localhost:6370:0 /tmp/slow-keys.txt 1


Beast progress bar courtesy of cheggaaa:
80973 / 5000000 [>---------------------------------------------] 1.62 % 8066/s 10m9s


Use with caution:
I've tested it with redis's of ~5 mil keys (~10 G mem) etc. If somehow the program terminates early (which hasn't happened for me), you can continue progress by:
load keys from redis_source into a set (SADD each key), set1 = source
repeat for keys on redis_target, set2 = target
redis-cli diff set1 set2 > remaining.txt
pass remaining.txt to redis-transfer



Todo:
Add an <opts> param for flags. We would then be able to pass '-o' which would DEL a key if it already exists (which will allow RESTORE to overwrite a key). etcetc.


pc!
