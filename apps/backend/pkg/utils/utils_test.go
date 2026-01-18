package utils

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenDID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"did:key:z6MkhaXgBZDnaWnM", 
			"did:key:z6MkhaXgBZDnaWnM", // Short enough, 
		},
		{
			"did:key:z6MkhaXgBZDnaWnMz6MkhaXgBZDnaWnMz6MkhaXgBZDnaWnM",
			"did:key:z6MkhaXg...gBZDnaWnM", // Should truncate
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, ShortenDID(tt.input))
	}
}

func TestNormalizeToken(t *testing.T) {
	raw := []byte("raw-token-data")
	encoded := base64.StdEncoding.EncodeToString(raw)

	t.Run("Decodes Base64 explicit", func(t *testing.T) {
		res, err := NormalizeToken(encoded, "base64")
		assert.NoError(t, err)
		assert.Equal(t, raw, res)
	})

	t.Run("Auto-detects Base64 in Raw", func(t *testing.T) {
		res, err := NormalizeToken(encoded, "raw")
		assert.NoError(t, err)
		assert.Equal(t, raw, res)
	})

	t.Run("Handles plain raw string", func(t *testing.T) {
		// If it's not base64, it should return as is
		res, err := NormalizeToken("not-base-64", "raw")
		assert.NoError(t, err)
		assert.Equal(t, []byte("not-base-64"), res)
	})
}

func TestValidateTokenFormat(t *testing.T) {
	t.Run("Rejects empty", func(t *testing.T) {
		err := ValidateTokenFormat("")
		assert.Error(t, err)
	})

	t.Run("Rejects too short", func(t *testing.T) {
		err := ValidateTokenFormat("abc")
		assert.Error(t, err)
	})

	t.Run("Accepts valid base64", func(t *testing.T) {
		valid := base64.StdEncoding.EncodeToString([]byte("valid-token-content"))
		err := ValidateTokenFormat(valid)
		assert.NoError(t, err)
	})
}