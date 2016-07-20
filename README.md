# redis-transfer

Concurrent bulk transfer of keys from one redis server to another.

## What it does?

Parallelizes the transfer of redis keys from a source to a target redis server. You can specify what keys to transfer via regex or from an input file which contains a list of keys. Redis-transfer will then spawn N go routines (threads) which pull (dump) keys from the source redis server & immediately push (restore) those keys to the target redis. As a bonus, gives you a pretty progress bar.


## Usage

```
 $ ./redis-transfer
usage: redis-transfer <from_redis> <to_redis> <regex_or_input_file> <concurrent_threads>

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

 * examples:
   redis-transfer localhost:6379 remotehost:6379:1:password "migrate:*" 50
   redis-transfer redis://localhost:6379 redis://user:password@remotehost?db=1 "migrate:*" 50
```



## Example benchmark output

Transferring 5 million keys using 40 concurrent threads in ~10 minutes.

```
47393 / 5000000 [=>---------------------------------------------------------------------------------------------------------------------------------] 0.95 % 8061/s 10m14
10 threads:
108367 / 5000000 [==>-------------------------------------------------------------------------------------------------------------------------------] 2.17 % 6457/s 12m37
1 thread:
4124 / 5000000 [>-----------------------------------------------------------------------------------------------------------------------------------] 0.08 % 2059/s 40m26
```



## Todo

- Add an <opts> param for flags.
- With <opts>, ability to pass '--replace' which would replace existing keys.



## Contributing

Pull requests welcome!

See [CONTRIBUTORS](CONTRIBUTORS.md).
