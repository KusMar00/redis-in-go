package main

import (
	"sync"
)

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}


var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET": set,
	"GET": get,
	"DEL": del,
	"EXISTS": exists,
	"HSET": hset,
	"HGET": hget,
	"HGETALL": hgetall,
	"HDEL": hdel,
	"HEXISTS": hexists,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func del(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'del' command"}
	}

	key := args[0].bulk

	SETsMu.Lock()
	delete(SETs, key)
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func exists (args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'exists' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	_, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "integer", num: 0}
	}

	return Value{typ: "integer", num: 1}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	key := args[0].bulk

	HSETsMu.RLock()
	hash, ok := HSETs[key]
	HSETsMu.RUnlock()


	if !ok {
		return Value{typ: "null"}
	}
	
	values := make([]Value, 0, len(hash)*2)

	for field, value := range hash {
		values = append(values, Value{typ: "bulk", bulk: field})
		values = append(values, Value{typ: "bulk", bulk: value})
	}

	return Value{typ: "array", array: values}
}

func hdel(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hdel' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.Lock()
	delete(HSETs[hash], key)
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hexists (args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hexists' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	_, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "integer", num: 0}
	}

	return Value{typ: "integer", num: 1}
}