# TCP Echo Server & Web Client

이 프로젝트는 Go 언어로 구현된 **TCP Echo Server**와, 이를 웹 브라우저에서 테스트할 수 있게 해주는 **WebSocket 기반 웹 클라이언트**로 구성되어 있습니다.

## 📂 프로젝트 구조 (Project Structure)

| 파일명 | 역할 | 설명 |
|---|---|---|
| **`main.go`** | **TCP Server** | 핵심 비즈니스 로직. 8080 포트에서 클라이언트의 연결을 받아 메시지를 그대로 반환(Echo)합니다. |
| **`web_server.go`** | **Bridge Server** | 8081 포트에서 웹 요청을 처리합니다. 브라우저(WebSocket)와 TCP 서버 간의 중계자 역할을 합니다. |
| **`index.html`** | **Web Client** | 사용자가 메시지를 입력하는 웹 인터페이스입니다. WebSocket을 통해 웹 서버와 통신합니다. |
| **`client.go`** | **CLI Client** | 터미널에서 서버를 테스트하기 위한 간단한 Go 클라이언트입니다. |

## 🏗️ 아키텍처 및 구현 이유 (Architecture & Rationale)

### 1. 왜 Go 언어인가?
*   **동시성(Concurrency)**: Go의 `Goroutine`은 매우 가벼워서, 수천 개의 동시 연결도 적은 리소스로 처리할 수 있습니다. (`go handleConnection(conn)`)
*   **성능**: C++에 준하는 성능을 내면서도 생산성이 높습니다.
*   **표준 라이브러리**: `net` 패키지가 강력하여 별도 프레임워크 없이도 견고한 TCP 서버 구현이 가능합니다.

### 2. 왜 웹 서버(Bridge)가 필요한가?
*   **브라우저의 제약**: 웹 브라우저는 보안상의 이유로 **Raw TCP 통신을 직접 할 수 없습니다**.
*   **해결책**: 브라우저가 통신할 수 있는 HTTP/WebSocket을 지원하는 중간 서버(`web_server.go`)를 두어, 브라우저의 요청을 TCP 패킷으로 변환해 TCP 서버에 전달하는 **브리지(Bridge) 패턴**을 사용했습니다.

### 3. 왜 WebSocket인가? (HTTP vs WebSocket)
초기에는 HTTP POST 방식을 사용했으나, 다음과 같은 이유로 **WebSocket**으로 최적화했습니다.

*   **오버헤드 감소**: HTTP는 메시지마다 연결(3-way handshake)을 맺고 끊어야 하지만, WebSocket은 한 번 연결하면 계속 유지됩니다.
*   **실시간성**: 양방향 통신이 가능하여 서버의 응답을 실시간으로 받을 수 있습니다.
*   **효율성**: `web_server.go`에서 클라이언트마다 하나의 TCP 연결을 유지(`Persistent Connection`)하도록 하여, 불필요한 재연결 비용을 없애고 성능을 극대화했습니다.

## 🚀 실행 방법

1. **TCP 서버 실행**
   ```bash
   go run main.go
   ```

2. **웹 서버 실행** (새 터미널)
   ```bash
   go run web_server.go
   ```

3. **테스트**
   *   브라우저: [http://localhost:8081](http://localhost:8081)
   *   CLI: `go run client.go`
