package main

import (
	mcp "github.com/ne275/gin-mcp"
	"github.com/gin-gonic/gin"
)

// newEngine 构造与 main 一致的路由与 gin-mcp 挂载；publicBase 为工具执行时的 BaseURL。
func newEngine(publicBase string) *gin.Engine {
	app := gin.New()
	app.Use(gin.Recovery())

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

	return app
}
