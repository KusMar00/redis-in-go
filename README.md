# redis-in-go

Simple Redis server clone written in Go. An in-memory key-value database. This code exposes a TCP server on port :6379 (standard for Redis), which can be interacted with through the [redis-cli](https://redis.io/docs/latest/develop/tools/cli/). Make sure no other redis server is already running on port :6379.

## Supported Commands

- `PING \[message\]` _Responds with PONG or the message given as argument._
- `SET key value` _Assigns the value to the key._
- `GET key` _Fetches the value associated with the key._
- `HSET key field value` _Assigns the value to the field under the key_
- `HGET key field` _Fetches the value associated with a field and a key_
- `HGETALL key` _Lists all fields and values associated with a key_

## How to run

Make sure you have `Go` installed. Then run

```bash
$ cd app
$ go run .
```

Example of interacting with the server using the `redis-cli`:

```bash
$ redis-cli
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> HSET users u1 markus
OK
127.0.0.1:6379> HGETALL users
1) "u1"
2) "markus
```
