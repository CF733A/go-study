package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	request, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	parts := strings.Fields(request)
	if len(parts) < 2 {
		fmt.Println("wrong request")
		return
	}

	method := parts[0]
	path := parts[1]

	if method == "GET" && path == "/" {
		html := `<!DOCTYPE html>
<html>
<head>
<title>Webserver</title>
</head>
<body>
hello world
</body>
</html>`
		response := fmt.Sprintf("HTTP/1.1 200 OK\nContent-Type: text/html\n\n%s\n", html)
		conn.Write([]byte(response))
	} else {
		response := "HTTP/1.1 404 Not Found\nContent-Type: text/plain\n\n404 Not Found\n"
		conn.Write([]byte(response))
	}
}

func main() {
	listener, _ := net.Listen("tcp", ":8080")
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}
}
