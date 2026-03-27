package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamableMCP_ListTools(t *testing.T) {
	ts := httptest.NewServer(newHTTPHandlerWithStreamableOpts("/mcp", "http://127.0.0.1:9", &mcp.StreamableHTTPOptions{JSONResponse: true}))
	defer ts.Close()

	ctx := context.Background()
	cl := mcp.NewClient(&mcp.Implementation{Name: "vkube-mcp-test", Version: "1"}, nil)
	session, err := cl.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:             ts.URL + "/mcp",
		DisableStandaloneSSE: true, // 避免初始化后常驻 GET/SSE，否则 httptest.Server.Close 会阻塞
	}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	tools, err := session.ListTools(ctx, nil)
	require.NoError(t, err)
	names := make([]string, 0, len(tools.Tools))
	for _, tool := range tools.Tools {
		names = append(names, tool.Name)
	}
	assert.Contains(t, names, "get_vkube_file_schema")
	assert.Contains(t, names, "deploy_vkube_file")
}

func TestStreamableMCP_CallDeployTool(t *testing.T) {
	ts := httptest.NewServer(newHTTPHandlerWithStreamableOpts("/mcp", "http://127.0.0.1:9", &mcp.StreamableHTTPOptions{JSONResponse: true}))
	defer ts.Close()

	ctx := context.Background()
	cl := mcp.NewClient(&mcp.Implementation{Name: "vkube-mcp-test", Version: "1"}, nil)
	session, err := cl.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:             ts.URL + "/mcp",
		DisableStandaloneSSE: true,
	}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "deploy_vkube_file",
		Arguments: map[string]any{
			"vkubeFilePath": "/tmp/test.yaml",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.Content)
	txt, ok := res.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, txt.Text, "vkube deploy -f /tmp/test.yaml")
}

func TestHealthz(t *testing.T) {
	ts := httptest.NewServer(newHTTPHandlerWithStreamableOpts("/mcp", "http://127.0.0.1:9", &mcp.StreamableHTTPOptions{JSONResponse: true}))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/healthz")
	require.NoError(t, err)
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
