package main

import (
    "fmt"
    "net"
    "os"
    "github.com/ThatGuyFromFinland/utils"
    "strings"
)

const (
    CONN_HOST = "123.123.123.123"
    CONN_PORT = "50500"
    CONN_TYPE = "tcp"
)

type Connections struct {
    Id []uint64
    Address []Address
}

type Address struct {
    IP string
    port string
}

func main() {
    // Listen for incoming connections.
    serverAddr := Address{IP: ip.GetIP(), port: CONN_PORT}
    l, err := net.Listen(CONN_TYPE, serverAddr.IP+":"+serverAddr.port)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Printf("Listening on %s:%s", serverAddr.IP, serverAddr.port)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)
    // Read the incoming connection into the buffer.
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }
    // Send a response back to person contacting us.
    // message = "Message received: " + message
    request := strings.SplitN(string(buf[:n]), ":", 2)
    action := request[0]
    body := ""
    if (len(request) > 1) {
        body = request[1]
        switch {
            case action == "JOIN":
                handleClientJoin(body, conn)
        }
    }
    fmt.Printf("Received action: %s", action)
    fmt.Printf("Received body: %s", body)
    // Close the connection when you're done with it.
    conn.Close()
}

func handleClientJoin(body string, conn net.Conn) {
    conn.Write([]byte(body))
}