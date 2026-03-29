package integration

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/goddhi/ucan-visualizer/internal/api"
	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/test/fixtures"
)

// --- Helper Functions ---

func createMultipartRequest(t *testing.T, url string, filename string, fileContent []byte) *http.Request {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)
	_, err = part.Write(fileContent)
	require.NoError(t, err)
	err = writer.Close()
	require.NoError(t, err)

	req, err := http.NewRequest("POST", url, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func setupServer() *httptest.Server {
	handler := api.SetupRouter()
	return httptest.NewServer(handler)
}


func TestHealthCheck(t *testing.T) {
	server := setupServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "healthy", result["status"])
}

func TestParseEndpoints(t *testing.T) {
	server := setupServer()
	defer server.Close()

	validToken, _ := fixtures.GenerateValidUCAN()
	complexChain, _ := fixtures.GenerateComplexChain()
	
	validTokenStr := base64.StdEncoding.EncodeToString(validToken)
	complexChainStr := base64.StdEncoding.EncodeToString(complexChain)

	t.Run("POST /parse/delegation (JSON)", func(t *testing.T) {
		payload := models.ParseRequest{Token: validTokenStr}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(server.URL+"/api/parse/delegation", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST /parse/delegation/file (File)", func(t *testing.T) {
		req := createMultipartRequest(t, server.URL+"/api/parse/delegation/file", "token.ucan", validToken)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST /parse/chain (JSON)", func(t *testing.T) {
		payload := models.ParseRequest{Token: complexChainStr}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(server.URL+"/api/parse/chain", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST /parse/chain/file (File)", func(t *testing.T) {
		req := createMultipartRequest(t, server.URL+"/api/parse/chain/file", "chain.ucan", complexChain)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST /parse/invocation (JSON)", func(t *testing.T) {
		payload := models.ParseRequest{Token: validTokenStr}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(server.URL+"/api/parse/invocation", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST /parse/invocation/file (File)", func(t *testing.T) {
		req := createMultipartRequest(t, server.URL+"/api/parse/invocation/file", "invoke.ucan", validToken)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestValidateEndpoints(t *testing.T) {
	server := setupServer()
	defer server.Close()

	complexChain, _ := fixtures.GenerateComplexChain()
	complexChainStr := base64.StdEncoding.EncodeToString(complexChain)

	t.Run("POST /validate/chain (JSON)", func(t *testing.T) {
		payload := models.ValidateRequest{Token: complexChainStr}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(server.URL+"/api/validate/chain", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		
		var result models.ValidationResult
		json.NewDecoder(resp.Body).Decode(&result)
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, result.Valid)
	})

	t.Run("POST /validate/chain/file (File)", func(t *testing.T) {
		req := createMultipartRequest(t, server.URL+"/api/validate/chain/file", "chain.ucan", complexChain)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		
		var result models.ValidationResult
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.True(t, result.Valid)
	})

	t.Run("Validate Expired Token", func(t *testing.T) {
		expiredToken, _ := fixtures.GenerateExpiredUCAN()
		expiredTokenStr := base64.StdEncoding.EncodeToString(expiredToken)

		payload := models.ValidateRequest{Token: expiredTokenStr}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(server.URL+"/api/validate/chain", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)

		var result models.ValidationResult
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.False(t, result.Valid)
	})
}

func TestGraphEndpoints(t *testing.T) {
	server := setupServer()
	defer server.Close()

	complexChain, _ := fixtures.GenerateComplexChain()
	complexChainStr := base64.StdEncoding.EncodeToString(complexChain)

	t.Run("POST /graph/delegation (JSON)", func(t *testing.T) {
		payload := models.GraphRequest{Token: complexChainStr}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(server.URL+"/api/graph/delegation", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		
		var result models.GraphResponse
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, result.Nodes)
		assert.NotEmpty(t, result.Edges)
	})

	t.Run("POST /graph/delegation/file (File)", func(t *testing.T) {
		req := createMultipartRequest(t, server.URL+"/api/graph/delegation/file", "graph.ucan", complexChain)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST /graph/invocation (JSON)", func(t *testing.T) {
		payload := models.GraphRequest{Token: complexChainStr}
		body, _ := json.Marshal(payload)
		resp, err := http.Post(server.URL+"/api/graph/invocation", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST /graph/invocation/file (File)", func(t *testing.T) {
		req := createMultipartRequest(t, server.URL+"/api/graph/invocation/file", "inv_graph.ucan", complexChain)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}