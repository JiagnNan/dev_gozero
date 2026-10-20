package main

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
	"time"
)

type SseHandler struct {
	clients map[chan string]bool
}

func NewSseHandler() *SseHandler {
	return &SseHandler{
		clients: make(map[chan string]bool),
	}
}

// Serve 处理 SSE 连接
func (h *SseHandler) Serve(w http.ResponseWriter, r *http.Request) {
	// 设置 SSE 必需的 HTTP 头
	// for versions > v1.8.1, no need to add 3 lines below
	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")

	// 为每个客户端创建一个 channel
	clientChan := make(chan string)
	h.clients[clientChan] = true

	// 客户端断开时清理
	defer func() {
		delete(h.clients, clientChan)
		close(clientChan)
	}()

	// 持续监听并推送事件
	for {
		select {
		case msg := <-clientChan:
			// 发送事件数据
			//收到消息时，格式化为 SSE 协议（data: 消息\n\n），并通过 Flush() 立即推送。
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			// 客户端断开连接
			return
		}
	}
}

// SimulateEvents 模拟周期性事件
func (h *SseHandler) SimulateEvents() {
	//使用 time.Ticker 每秒生成一个事件（当前时间）
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		message := fmt.Sprintf("Server time: %s", time.Now().Format(time.RFC3339))
		// 广播给所有客户端
		for clientChan := range h.clients {
			select {
			case clientChan <- message:
			default:
				// 跳过阻塞的 channel
			}
		}
	}
}

func main() {
	// 创建 go-zero REST 服务，集成静态文件服务
	server := rest.MustNewServer(rest.RestConf{
		Host: "0.0.0.0",
		Port: 8080,
	}, rest.WithFileServer("/static", http.Dir("static"))) //配置静态文件服务，映射 /static 到本地 static 目录
	defer server.Stop()

	// 初始化 SSE 处理
	sseHandler := NewSseHandler()

	// 注册 SSE 路由
	//// for go-zero versions <= v1.8.1
	//server.AddRoute(rest.Route{
	//	Method:  http.MethodGet,
	//	Path:    "/sse",
	//	Handler: sseHandler.Serve,
	//}, rest.WithTimeout(0))

	// for go-zero versions > v1.8.1
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/sse",
		Handler: sseHandler.Serve,
	}, rest.WithSSE())

	// 在单独的 goroutine 中模拟事件
	go sseHandler.SimulateEvents()

	logx.Info("Server starting on :8080")
	server.Start()
}
