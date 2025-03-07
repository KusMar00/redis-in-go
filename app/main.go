package main

import (
	"fmt"
	"net"
)


func main() {
	fmt.Println("Listening on port :6379")

	// create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	// listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting:", err)
		return
	}

	defer conn.Close() // close connection once finished

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		_ = value

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}