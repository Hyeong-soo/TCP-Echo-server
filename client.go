package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 1. 서버에 연결합니다.
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to server. Type your message and press Enter (type 'quit' to exit).")

	// 입력을 받기 위한 Reader
	inputReader := bufio.NewReader(os.Stdin)
	// 서버 응답을 받기 위한 Reader
	serverReader := bufio.NewReader(conn)

	for {
		// 2. 사용자 입력 받기
		fmt.Print("Enter message: ")
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "quit" || text == "exit" {
			fmt.Println("Exiting...")
			break
		}

		// 3. 서버로 메시지 전송 (줄바꿈 추가)
		_, err = fmt.Fprintf(conn, "%s\n", text)
		if err != nil {
			fmt.Println("Error writing:", err)
			break
		}

		// 4. 서버로부터 응답 수신
		response, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}

		fmt.Printf("Server response: %s", response)
	}
}
