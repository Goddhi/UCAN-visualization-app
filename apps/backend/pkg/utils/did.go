package utils

import (
	"fmt"
	"strings"
)

func ShortenDID(did string) string {
	if len(did) <= 25 {
		return did
	}

	parts := strings.Split(did, ":")
	if len(parts) < 3 {
		return did
	}

	scheme := parts[0]     
	method := parts[1]     
	identifier := parts[2] 

	if len(identifier) > 12 {
		return fmt.Sprintf("%s:%s:%s...%s",
			scheme,
			method,
			identifier[:6],
			identifier[len(identifier)-3:],
		)
	}

	return did
}

func ValidateDID(did string) bool {
	parts := strings.Split(did, ":")
	if len(parts) < 3 {
		return false
	}

	if parts[0] != "did" {
		return false
	}

	// Method must be non-empty
	if parts[1] == "" {
		return false
	}

	// Method-specific ID must be non-empty
	if parts[2] == "" {
		return false
	}

	return true
}