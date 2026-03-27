package main

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer_RealTCPPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	addr := ln.Addr().String()
	baseURL := "http://" + addr
	h := newHTTPHandlerWithStreamableOpts("/mcp", baseURL, &mcp.StreamableHTTPOptions{JSONResponse: true})

	srv := &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.Serve(ln) }()
	defer func() {
		_ = srv.Close()
		select {
		case <-errCh:
		case <-time.After(2 * time.Second):
		}
	}()

	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("healthz", func(t *testing.T) {
		res, err := client.Get(baseURL + "/healthz")
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("mcp_post_initialize", func(t *testing.T) {
		body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}`
		req, err := http.NewRequest(http.MethodPost, baseURL+"/mcp", strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		res, err := client.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		assert.GreaterOrEqual(t, res.StatusCode, 200)
		assert.Less(t, res.StatusCode, 300, "body=%s", string(b))
	})
}

func TestHTTPServer_httptestMCPPath(t *testing.T) {
	ts := httptestNewServerWithMCP(t)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/healthz")
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func httptestNewServerWithMCP(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(newHTTPHandlerWithStreamableOpts("/mcp", "http://127.0.0.1:1", &mcp.StreamableHTTPOptions{JSONResponse: true}))
}
