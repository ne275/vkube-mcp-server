package main

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHTTPServer_RealTCPPort 在真实 TCP 端口上起服务并发 HTTP 请求（与 httptest 内存复现不同）。
func TestHTTPServer_RealTCPPort(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	addr := ln.Addr().String()
	baseURL := "http://" + addr
	app := newEngine(baseURL)

	srv := &http.Server{
		Handler:           app,
		ReadHeaderTimeout: 5 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(ln)
	}()
	defer func() {
		_ = srv.Close()
		select {
		case <-errCh:
		case <-time.After(2 * time.Second):
		}
	}()

	client := &http.Client{Timeout: 5 * time.Second}

	t.Run("vkubefile", func(t *testing.T) {
		res, err := client.Get(baseURL + "/api/v1/mcp/vkubefile")
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		b, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Contains(t, string(b), "jsonrpc")
		assert.Contains(t, string(b), "VKubeFile")
	})

	t.Run("deployCommand", func(t *testing.T) {
		res, err := client.Get(baseURL + "/api/v1/mcp/deployCommand")
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		b, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Contains(t, string(b), "vkube deploy -f")
	})
}

// TestHTTPServer_httptestServer 使用 httptest 在随机端口监听并发请求。
func TestHTTPServer_httptestServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ts := httptest.NewUnstartedServer(newEngine("http://127.0.0.1:1"))
	ts.Start()
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/v1/mcp/vkubefile")
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
