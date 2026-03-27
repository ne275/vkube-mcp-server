// vkube-mcp：独立 MCP 服务，使用官方 go-sdk 的 Streamable HTTP（与 Cursor 兼容）。
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	listen := getenv("VKUBE_MCP_LISTEN", ":3100")
	mcpPath := getenv("VKUBE_MCP_HTTP_PATH", "/mcp")
	publicHint := getenv("VKUBE_MCP_PUBLIC_BASE_URL", defaultPublicBaseURL(listen))

	h := newHTTPHandler(mcpPath, publicHint)
	base := strings.TrimSuffix(publicHint, "/")
	log.Printf("vkube-mcp: MCP Streamable HTTP at path %q — set Cursor mcp.json url to: %s%s", mcpPath, base, mcpPath)

	if err := listenAndServe(listen, h); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}

func defaultPublicBaseURL(listen string) string {
	listen = strings.TrimSpace(listen)
	if strings.HasPrefix(listen, ":") {
		return "http://127.0.0.1" + listen
	}
	host, port, err := net.SplitHostPort(listen)
	if err != nil {
		return "http://127.0.0.1:3100"
	}
	if host == "" || host == "0.0.0.0" || host == "[::]" || host == "::" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("http://%s:%s", host, port)
}
