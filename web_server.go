package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// 모든 오리진 허용 (개발 편의상)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Web Server running on http://localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("Error starting web server:", err)
		os.Exit(1)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. HTTP 요청을 WebSocket 연결로 업그레이드
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	// 2. 백엔드 TCP 서버와 연결 (클라이언트당 1개의 TCP 연결 유지)
	tcpConn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("TCP Dial error:", err)
		ws.WriteMessage(websocket.TextMessage, []byte("Error connecting to TCP server"))
		return
	}
	defer tcpConn.Close()

	// 3. 고루틴을 사용하여 양방향 통신 처리
	// 채널을 사용하여 에러 발생 시 두 루프를 모두 종료하도록 함
	done := make(chan struct{})

	// Goroutine: TCP 서버 -> WebSocket (수신)
	go func() {
		defer close(done)
		reader := bufio.NewReader(tcpConn)
		for {
			// TCP 서버로부터 메시지 읽기
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("TCP Read error:", err)
				return
			}

			// WebSocket 클라이언트에게 전송
			err = ws.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("WS Write error:", err)
				return
			}
		}
	}()

	// Main Loop: WebSocket -> TCP 서버 (송신)
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("WS Read error:", err)
			break
		}

		// TCP 서버로 전송 (줄바꿈 추가)
		_, err = tcpConn.Write(append(message, '\n'))
		if err != nil {
			fmt.Println("TCP Write error:", err)
			break
		}
	}

	// 종료 시그널 대기 (TCP 수신 고루틴이 끝날 때까지)
	<-done
}
