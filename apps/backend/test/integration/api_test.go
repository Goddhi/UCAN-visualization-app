package integration

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/goddhi/ucan-visualizer/internal/api"
	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/test/fixtures"
)

func TestParseEndpoint(t *testing.T) {
	// Setup test server
	handler := api.SetupRouter()
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Run("Parse valid delegation", func(t *testing.T) {
		// Generate a valid UCAN
		tokenBytes, err := fixtures.GenerateValidUCAN()
		require.NoError(t, err)

		// Encode as base64 for JSON transport
		tokenStr := base64.StdEncoding.EncodeToString(tokenBytes)

		// Prepare request
		payload := models.ParseRequest{
			Token: tokenStr,
		}
		body, _ := json.Marshal(payload)

		// Make request
		resp, err := http.Post(
			server.URL+"/api/parse/delegation",
			"application/json",
			bytes.NewBuffer(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Check response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		var result models.DelegationResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		// Verify parsed data
		assert.NotEmpty(t, result.Issuer)
		assert.NotEmpty(t, result.Audience)
		assert.NotEmpty(t, result.CID)
		assert.Greater(t, len(result.Capabilities), 0)

		t.Logf("✅ Parsed delegation: %s", result.CID)
		t.Logf("   Issuer: %s", result.Issuer)
		t.Logf("   Audience: %s", result.Audience)
		t.Logf("   Capabilities: %d", len(result.Capabilities))
	})

	t.Run("Parse with invalid token", func(t *testing.T) {
		payload := models.ParseRequest{
			Token: "invalid-token",
		}
		body, _ := json.Marshal(payload)

		resp, err := http.Post(
			server.URL+"/api/parse/delegation",
			"application/json",
			bytes.NewBuffer(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Parse with empty token", func(t *testing.T) {
		payload := models.ParseRequest{
			Token: "",
		}
		body, _ := json.Marshal(payload)

		resp, err := http.Post(
			server.URL+"/api/parse/delegation",
			"application/json",
			bytes.NewBuffer(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestValidateEndpoint(t *testing.T) {
	handler := api.SetupRouter()
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Run("Validate valid delegation", func(t *testing.T) {
		// Generate a valid UCAN
		tokenBytes, err := fixtures.GenerateValidUCAN()
		require.NoError(t, err)

		tokenStr := base64.StdEncoding.EncodeToString(tokenBytes)

		payload := models.ValidateRequest{
			Token: tokenStr,
		}
		body, _ := json.Marshal(payload)

		resp, err := http.Post(
			server.URL+"/api/validate/chain",
			"application/json",
			bytes.NewBuffer(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result models.ValidationResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		// Check validation result
		t.Logf("✅ Validation result: Valid=%v", result.Valid)
		t.Logf("   Total links: %d", result.Summary.TotalLinks)
		t.Logf("   Valid links: %d", result.Summary.ValidLinks)
		t.Logf("   Invalid links: %d", result.Summary.InvalidLinks)

		if len(result.Chain) > 0 {
			t.Logf("   First link issues: %d", len(result.Chain[0].Issues))
			for _, issue := range result.Chain[0].Issues {
				t.Logf("     - [%s] %s", issue.Severity, issue.Message)
			}
		}
	})

	t.Run("Validate expired delegation", func(t *testing.T) {
		tokenBytes, err := fixtures.GenerateExpiredUCAN()
		require.NoError(t, err)

		tokenStr := base64.StdEncoding.EncodeToString(tokenBytes)

		payload := models.ValidateRequest{
			Token: tokenStr,
		}
		body, _ := json.Marshal(payload)

		resp, err := http.Post(
			server.URL+"/api/validate/chain",
			"application/json",
			bytes.NewBuffer(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result models.ValidationResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		// Should be invalid due to expiration
		assert.False(t, result.Valid)
		assert.NotNil(t, result.RootCause)
		assert.Equal(t, "expired", result.RootCause.Type)

		t.Logf("✅ Correctly detected expired UCAN")
	})
}

func TestGraphEndpoint(t *testing.T) {
	handler := api.SetupRouter()
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Run("Generate graph from valid delegation", func(t *testing.T) {
		tokenBytes, err := fixtures.GenerateValidUCAN()
		require.NoError(t, err)

		tokenStr := base64.StdEncoding.EncodeToString(tokenBytes)

		payload := models.GraphRequest{
			Token: tokenStr,
		}
		body, _ := json.Marshal(payload)

		resp, err := http.Post(
			server.URL+"/api/graph/delegation",
			"application/json",
			bytes.NewBuffer(body),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result models.GraphResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		// Check graph structure
		assert.GreaterOrEqual(t, len(result.Nodes), 2, "Should have at least 2 nodes (issuer and audience)")
		assert.Greater(t, len(result.Edges), 0, "Should have at least 1 edge")

		t.Logf("✅ Generated graph:")
		t.Logf("   Nodes: %d", len(result.Nodes))
		t.Logf("   Edges: %d", len(result.Edges))

		for i, node := range result.Nodes {
			t.Logf("   Node %d: %s (%s)", i+1, node.Label, node.Type)
		}

		for i, edge := range result.Edges {
			t.Logf("   Edge %d: %s -> %s (%s)", i+1, edge.Source, edge.Target, edge.Label)
		}
	})
}

func TestHealthEndpoint(t *testing.T) {
	handler := api.SetupRouter()
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "healthy", result["status"])
	t.Logf("✅ Health check passed")
}