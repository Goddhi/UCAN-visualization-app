package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/storacha/go-ucanto/core/delegation"
	"github.com/storacha/go-ucanto/principal/ed25519/signer"
	"github.com/storacha/go-ucanto/ucan"
)

func main() {
	// Create issuer (Alice)
	alice, err := signer.Generate()
	if err != nil {
		log.Fatalf("Failed to generate issuer: %v", err)
	}

	// Create audience (Bob)
	bob, err := signer.Generate()
	if err != nil {
		log.Fatalf("Failed to generate audience: %v", err)
	}

	// Create a delegation
	del, err := delegation.Delegate(
		alice,
		bob,
		[]ucan.Capability[ucan.NoCaveats]{
			ucan.NewCapability(
				"store/add",
				"storage:alice/*",
				ucan.NoCaveats{},
			),
		},
		delegation.WithExpiration(int(time.Now().Add(24*time.Hour).Unix())),
		delegation.WithNotBefore(int(time.Now().Unix())),
	)
	if err != nil {
		log.Fatalf("Failed to create delegation: %v", err)
	}

	// Archive as CAR bytes
	archive := del.Archive()
	carBytes, err := io.ReadAll(archive)
	if err != nil {
		log.Fatalf("Failed to read archive: %v", err)
	}

	// Encode as base64
	tokenStr := base64.StdEncoding.EncodeToString(carBytes)

	fmt.Println("Generated UCAN Token:")
	fmt.Println("=====================")
	fmt.Println(tokenStr)
	fmt.Println()
	fmt.Println("Copy this token and use it in your tests")
	fmt.Println()
	fmt.Printf("Test command:\n")
	fmt.Printf(`curl -X POST http://localhost:8080/api/parse/delegation \
  -H "Content-Type: application/json" \
  -d '{"token":"%s"}' | jq
`, tokenStr)
}