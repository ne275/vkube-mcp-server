package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newMCPServer 创建带工具的 MCP Server（会话模式，支持 Cursor 的 Streamable HTTP + SSE 回退）。
func newMCPServer() *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "vkube-mcp",
		Version: "1.0.0",
	}, nil)
	registerTools(s)
	return s
}

func streamableHTTPOptions() *mcp.StreamableHTTPOptions {
	o := &mcp.StreamableHTTPOptions{}
	if truthy(os.Getenv("VKUBE_MCP_DISABLE_LOCALHOST_PROTECTION")) {
		o.DisableLocalhostProtection = true
	}
	// 默认使用 JSON 响应而非长驻 SSE，便于反向代理与 Cursor 以 POST 为主的 Streamable HTTP 流程。
	// 若需旧版纯 SSE，设置环境变量 VKUBE_MCP_USE_SSE=1
	if !truthy(os.Getenv("VKUBE_MCP_USE_SSE")) {
		o.JSONResponse = true
	}
	return o
}

func truthy(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "1" || s == "true" || s == "yes"
}

// newHTTPHandler 返回根路由：健康检查、说明页、MCP Streamable HTTP 端点。
// mcpPath 应为 "/mcp" 形式（无尾随斜杠也可由 mountMCP 同时注册 /mcp/）。
func newHTTPHandler(mcpPath, publicHint string) http.Handler {
	return newHTTPHandlerWithStreamableOpts(mcpPath, publicHint, streamableHTTPOptions())
}

// newHTTPHandlerWithStreamableOpts 与 newHTTPHandler 相同，但可注入 Streamable 选项（测试中使用 JSONResponse 避免长连接挂住 httptest）。
func newHTTPHandlerWithStreamableOpts(mcpPath, publicHint string, streamOpts *mcp.StreamableHTTPOptions) http.Handler {
	mux := http.NewServeMux()
	mcpServer := newMCPServer()
	stream := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return mcpServer
	}, streamOpts)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "vkube-mcp\n\nConfigure Cursor mcp.json \"url\" to the MCP endpoint (must include path), e.g.:\n  %s%s\n\nHealth: GET /healthz\n", strings.TrimSuffix(publicHint, "/"), mcpPath)
	})

	mountMCP(mux, mcpPath, stream)

	return mux
}

func mountMCP(mux *http.ServeMux, path string, h http.Handler) {
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		path = "/mcp"
	}
	mux.Handle(path, h)
	mux.Handle(path+"/", h)
}

func listenAndServe(addr string, h http.Handler) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      0,
	}
	log.Printf("vkube-mcp listening on %s", addr)
	return srv.ListenAndServe()
}
