// vkube-mcp：与 v-kube-service 无关的独立 MCP（SSE）服务，仅依赖 gin 与 gin-mcp。
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	mcp "github.com/ne275/gin-mcp"
	"github.com/gin-gonic/gin"
)

func main() {
	listen := getenv("VKUBE_MCP_LISTEN", ":3100")
	publicBase := getenv("VKUBE_MCP_PUBLIC_BASE_URL", defaultPublicBaseURL(listen))

	if strings.EqualFold(os.Getenv("GIN_MODE"), "release") {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()

	v1 := app.Group("/api/v1")
	m := v1.Group("/mcp")
	{
		m.GET("/vkubefile", handleGetVKubeFileSchema)
		m.GET("/deployCommand", handleDeployCommand)
	}

	srv := mcp.New(app, &mcp.Config{
		Name:        "vkube-mcp",
		Description: "VKubeFile schema and deploy command helpers (standalone MCP server)",
		BaseURL:     publicBase,
	})
	srv.RegisterSchema("GET", "/api/v1/mcp/vkubefile", nil, nil)
	srv.RegisterSchema("GET", "/api/v1/mcp/deployCommand", nil, nil)
	srv.Mount("/api/v1/mcp")

	log.Printf("vkube-mcp listen=%s public_base_url=%s", listen, publicBase)
	if err := app.Run(listen); err != nil {
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
	if host == "" || host == "0.0.0.0" || host == "[::]" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("http://%s:%s", host, port)
}
