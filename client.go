package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 1. 서버에 연결합니다.
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 2. 서버로 메시지를 전송합니다.
	message := "Hello, TCP Server!\n"
	fmt.Printf("Sending: %s", message)
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}

	// 3. 서버로부터 응답을 수신합니다.
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	fmt.Printf("Received: %s", response)
}
