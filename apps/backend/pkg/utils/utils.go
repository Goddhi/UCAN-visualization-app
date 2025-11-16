package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

// NormalizeToken converts token from various formats to bytes for CAR parsing
func NormalizeToken(token, format string) ([]byte, error) {
	token = strings.TrimSpace(token)
	
	switch strings.ToLower(format) {
	case "base64":
		// Decode base64 to get CAR bytes
		decoded, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return nil, fmt.Errorf("invalid base64 token: %w", err)
		}
		return decoded, nil
	case "raw", "":
		// For CAR format, we need to decode base64 by default
		// since our token generator outputs base64-encoded CAR data
		if isBase64(token) {
			decoded, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				// If base64 decode fails, return as raw bytes
				return []byte(token), nil
			}
			return decoded, nil
		}
		return []byte(token), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// ReadUploadedFile reads the contents of an uploaded file
func ReadUploadedFile(file multipart.File, header *multipart.FileHeader) ([]byte, error) {
	// Check file size (limit to 10MB)
	if header.Size > 10*1024*1024 {
		return nil, fmt.Errorf("file too large: %d bytes (max 10MB)", header.Size)
	}

	// Read file contents
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	// For CAR format files, check if content is base64 and decode if needed
	content := strings.TrimSpace(string(data))
	if isBase64(content) {
		decoded, err := base64.StdEncoding.DecodeString(content)
		if err == nil {
			return decoded, nil
		}
	}

	return data, nil
}

// ShortenDID creates a readable short version of a DID for display
func ShortenDID(did string) string {
	if len(did) <= 20 {
		return did
	}
	
	// Extract method and first/last parts
	parts := strings.Split(did, ":")
	if len(parts) < 3 {
		return did[:20] + "..."
	}
	
	method := parts[1]
	identifier := parts[2]
	
	if len(identifier) > 16 {
		return fmt.Sprintf("did:%s:%s...%s", 
			method, 
			identifier[:8], 
			identifier[len(identifier)-8:])
	}
	
	return did
}

// ValidateTokenFormat checks if the provided token appears to be valid
func ValidateTokenFormat(token string) error {
	token = strings.TrimSpace(token)
	
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}
	
	// Basic length check
	if len(token) < 10 {
		return fmt.Errorf("token too short")
	}
	
	// For CAR format, tokens are typically base64 encoded
	if isBase64(token) {
		// Try to decode to check if it's valid base64
		_, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return fmt.Errorf("invalid base64 encoding: %w", err)
		}
		return nil
	}
	
	// If not base64, allow raw bytes
	return nil
}

// isBase64 checks if a string is valid base64
func isBase64(s string) bool {
	// Check if string looks like base64 (contains only base64 characters)
	if len(s) == 0 {
		return false
	}
	
	// Base64 strings should be divisible by 4 (with padding) or close to it
	// and contain only valid base64 characters
	for _, char := range s {
		if !((char >= 'A' && char <= 'Z') || 
			 (char >= 'a' && char <= 'z') || 
			 (char >= '0' && char <= '9') || 
			 char == '+' || char == '/' || char == '=') {
			return false
		}
	}
	
	// Try to decode to confirm it's valid base64
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// ExtractFileExtension gets the file extension from filename
func ExtractFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
	}
	return strings.ToLower(parts[len(parts)-1])
}

// IsValidUCANFile checks if uploaded file appears to contain UCAN data
func IsValidUCANFile(data []byte, filename string) error {
	// Check file extension
	ext := ExtractFileExtension(filename)
	validExts := map[string]bool{
		"txt": true, "ucan": true, "car": true, 
		"token": true, "json": true, "":  true,
	}
	
	if !validExts[ext] {
		return fmt.Errorf("unsupported file type: %s", ext)
	}
	
	// Basic content validation
	content := string(data)
	if err := ValidateTokenFormat(content); err != nil {
		return fmt.Errorf("file does not contain valid UCAN data: %w", err)
	}
	
	return nil
}