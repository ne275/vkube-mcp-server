package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testEngine(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	return newEngine("http://127.0.0.1:9")
}

func TestHandleGetVKubeFileSchema_HTTP(t *testing.T) {
	app := testEngine(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/mcp/vkubefile", nil)
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			Content []struct {
				Type     string `json:"type"`
				Text     string `json:"text"`
				MimeType string `json:"mimeType"`
			} `json:"content"`
		} `json:"result"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "2.0", body.Jsonrpc)
	require.Len(t, body.Result.Content, 2)
	assert.Equal(t, "text", body.Result.Content[0].Type)
	assert.Contains(t, body.Result.Content[0].Text, "VKubeFile")
	assert.Equal(t, "resource", body.Result.Content[1].Type)
	assert.Equal(t, "application/json", body.Result.Content[1].MimeType)
	assert.Contains(t, body.Result.Content[1].Text, `"Kind"`)
	assert.Contains(t, body.Result.Content[1].Text, `"containers"`)
}

func TestHandleDeployCommand_HTTP_DefaultPath(t *testing.T) {
	app := testEngine(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/mcp/deployCommand", nil)
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body struct {
		Result struct {
			Command     string `json:"command"`
			DefaultPath string `json:"defaultPath"`
		} `json:"result"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "/path/to/vkubefile.yaml", body.Result.DefaultPath)
	assert.Equal(t, "vkube deploy -f /path/to/vkubefile.yaml", body.Result.Command)
}

func TestHandleDeployCommand_HTTP_CustomPathInBody(t *testing.T) {
	app := testEngine(t)
	w := httptest.NewRecorder()
	body := strings.NewReader(`{"vkubeFilePath":"/home/u/app/vkube.yaml"}`)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/mcp/deployCommand", body)
	req.Header.Set("Content-Type", "application/json")
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Result struct {
			Command     string `json:"command"`
			DefaultPath string `json:"defaultPath"`
		} `json:"result"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "/home/u/app/vkube.yaml", resp.Result.DefaultPath)
	assert.Equal(t, "vkube deploy -f /home/u/app/vkube.yaml", resp.Result.Command)
}

func TestToJSON(t *testing.T) {
	out := toJSON(map[string]int{"a": 1})
	assert.Contains(t, out, `"a"`)
	assert.Contains(t, out, `1`)
}
