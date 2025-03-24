package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
    fmt.Println("Listening on port :6379")

    // create a new server
    l, err := net.Listen("tcp", ":6379")
    if err != nil {
        fmt.Println("Error listening:", err)
        return
    }
    defer l.Close()

    // load AOF
    aof, err := NewAof("database.aof")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer aof.Close()

    // Replay AOF commands
    aof.Read(func(value Value) {
        command := strings.ToUpper(value.array[0].bulk)
        args := value.array[1:]

        handler, ok := Handlers[command]
        if !ok {
            fmt.Println("Invalid command: ", command)
            return
        }

        handler(args)
    })

    // Accept and handle multiple client connections
    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting:", err)
            continue
        }

        fmt.Println("Received connection", conn.RemoteAddr())

        // Handle the client in a separate goroutine
        go handleClient(conn, aof)
    }
}

func handleClient(conn net.Conn, aof *Aof) {
    defer conn.Close() // Ensure the connection is closed when the function exits

    resp := NewResp(conn)
    writer := NewWriter(conn)

    for {
        // Read the client's request
        value, err := resp.Read()
        if err != nil {
            fmt.Println("Error reading from client:", err)
            return
        }

        if value.typ != "array" {
            fmt.Println("Invalid request, expected array")
            writer.Write(Value{typ: "error", str: "ERR invalid request, expected array"})
            continue
        }

        if len(value.array) == 0 {
            fmt.Println("Invalid request, expected array length > 0")
            writer.Write(Value{typ: "error", str: "ERR invalid request, expected array length > 0"})
            continue
        }

        // Parse the command and arguments
        command := strings.ToUpper(value.array[0].bulk)
        args := value.array[1:]

        // Find the appropriate handler
        handler, ok := Handlers[command]
        if !ok {
            fmt.Println("Invalid command:", command)
            writer.Write(Value{typ: "error", str: "ERR unknown command"})
            continue
        }

        // Write to AOF for certain commands
        if command == "SET" || command == "HSET" || command == "DEL" || command == "HDEL" {
            aof.Write(value)
        }

        // Execute the command and write the result back to the client
        result := handler(args)
        writer.Write(result)
    }
}