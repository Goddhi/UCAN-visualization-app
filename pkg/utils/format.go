package utils

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
)

type TokenFormat string

const (
	FormatBase64 TokenFormat = "base64"
	FormatBinary TokenFormat = "binary"
	FormatHex    TokenFormat = "hex"
	FormatCID    TokenFormat = "cid"
	FormatUnknown TokenFormat = "unknown"
)

func DetectTokenFormat(input string) TokenFormat {
	if len(input) == 0 {
		return FormatUnknown
	}
	// Check for CID format
	if strings.HasPrefix(input, "bafy") || 
	   strings.HasPrefix(input, "bafk") || 
	   strings.HasPrefix(input, "Qm") {
		return FormatCID
	}
	// Check for hex format
	if strings.HasPrefix(input, "0x") || isHexString(strings.TrimPrefix(input, "0x")) {
		return FormatHex
	}
	// Check if it's valid base64
	if isBase64(input) {
		return FormatBase64
	}
	// If all else fails, assume binary 
	return FormatBinary
}
// NormalizeToken converts any format to raw CAR bytes
func NormalizeToken(input string, format string) ([]byte, error) {
	// Determine format
	var detectedFormat TokenFormat
	if format != "" && format != "auto" {
		detectedFormat = TokenFormat(format)
	} else {
		detectedFormat = DetectTokenFormat(input)
	}
	// Convert to bytes based on format
	switch detectedFormat {
	case FormatBase64:
		return base64.StdEncoding.DecodeString(input)
		
	case FormatHex:
		hexStr := strings.TrimPrefix(input, "0x")
		return hex.DecodeString(hexStr)
		
	case FormatCID:
		return nil, errors.New("CID format not yet supported - please provide the CAR data directly")
	case FormatBinary:
		// Already binary (from file upload)
		return []byte(input), nil
	default:
		return nil, errors.New("unknown token format")
	}
}

// isBase64 checks if a string is valid base64
func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// isHexString checks if a string contains only hex characters
func isHexString(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}