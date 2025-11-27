package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

// main 함수는 프로그램의 진입점입니다.
func main() {
	// 1. TCP 서버를 시작합니다.
	// net.Listen은 지정된 네트워크(tcp)와 주소(:8080)에서 연결을 기다립니다.
	// ":8080"은 모든 IP 주소의 8080 포트를 의미합니다.
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		// 서버 시작에 실패하면 에러를 출력하고 프로그램을 종료합니다.
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	// defer는 함수가 종료되기 직전에 실행됨을 보장합니다.
	// 서버가 종료될 때 리스너를 닫아 리소스를 정리합니다.
	defer listener.Close()

	fmt.Println("TCP Echo Server is running on port 8080...")

	for {
		// 2. 클라이언트의 연결 요청을 수락(Accept)합니다.
		// Accept는 클라이언트가 연결될 때까지 대기(Block)합니다.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue // 에러가 발생해도 서버를 멈추지 않고 다음 연결을 기다립니다.
		}

		// 3. 연결된 클라이언트를 처리합니다.
		// 'go' 키워드는 새로운 고루틴(Goroutine)을 생성하여 함수를 비동기적으로 실행합니다.
		// 이를 통해 여러 클라이언트를 동시에 처리할 수 있습니다 (동시성).
		go handleConnection(conn)
	}
}

// handleConnection은 클라이언트와의 연결을 처리하는 함수입니다.
func handleConnection(conn net.Conn) {
	// 함수 종료 시 연결을 닫습니다.
	defer conn.Close()

	// 클라이언트의 주소를 출력합니다.
	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	// 4. 클라이언트로부터 데이터를 읽고 그대로 돌려줍니다 (Echo).
	// bufio.NewReader를 사용하여 효율적으로 데이터를 읽습니다.
	reader := bufio.NewReader(conn)

	for {
		// 줄바꿈 문자('\n')를 만날 때까지 데이터를 읽습니다.
		// 텔넷(telnet)이나 넷캣(nc)은 엔터를 치면 '\n'을 보냅니다.
		message, err := reader.ReadString('\n')
		if err != nil {
			// 클라이언트가 연결을 끊거나(EOF) 에러가 발생하면 루프를 종료합니다.
			if err != io.EOF {
				fmt.Println("Error reading:", err)
			}
			break
		}

		// 수신한 메시지를 서버 콘솔에 출력합니다.
		fmt.Print("Received: ", message)

		// 5. 수신한 데이터를 그대로 클라이언트에게 다시 전송합니다 (Write).
		// 문자열을 바이트 슬라이스로 변환하여 전송합니다.
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing:", err)
			break
		}
	}

	fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())
}
