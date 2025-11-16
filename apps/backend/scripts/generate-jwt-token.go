package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/storacha/go-ucanto/core/delegation"
	"github.com/storacha/go-ucanto/principal/ed25519/signer"
	"github.com/storacha/go-ucanto/ucan"
)

func main() {
	// Generate signers using go-ucanto
	alice, err := signer.Generate()
	if err != nil {
		log.Fatal("Failed to generate alice:", err)
	}

	bob, err := signer.Generate()
	if err != nil {
		log.Fatal("Failed to generate bob:", err)
	}

	// Create capabilities using go-ucanto
	capabilities := []ucan.Capability[ucan.NoCaveats]{
		// Storage capabilities
		ucan.NewCapability("store/add", "storage:*", ucan.NoCaveats{}),
		ucan.NewCapability("store/get", "storage:*", ucan.NoCaveats{}),
		ucan.NewCapability("store/remove", "storage:*", ucan.NoCaveats{}),
		
		// Space capabilities  
		ucan.NewCapability("space/blob/add", "space:*", ucan.NoCaveats{}),
		ucan.NewCapability("space/index/add", "space:*", ucan.NoCaveats{}),
		
		// Upload capabilities
		ucan.NewCapability("upload/add", "upload:*", ucan.NoCaveats{}),
	}

	// Create delegation using go-ucanto
	del, err := delegation.Delegate(
		alice, // issuer
		bob,   // audience  
		capabilities,
		delegation.WithExpiration(int(time.Now().Add(24*time.Hour).Unix())),
		delegation.WithNotBefore(int(time.Now().Unix())),
	)
	if err != nil {
		log.Fatal("Failed to create delegation:", err)
	}

	// Export delegation to CAR format using go-ucanto
	archive := del.Archive() // This returns io.Reader
	
	// Read the archive data
	carData, err := io.ReadAll(archive)
	if err != nil {
		log.Fatal("Failed to read CAR archive:", err)
	}

	// Encode to base64 for transport
	token := base64.StdEncoding.EncodeToString(carData)

	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--token-only":
			fmt.Print(token)
			return
		case "--save":
			filename := "delegation.ucan"
			if len(os.Args) > 2 {
				filename = os.Args[2]
			}
			err := os.WriteFile(filename, carData, 0644)
			if err != nil {
				log.Fatal("Failed to save delegation:", err)
			}
			fmt.Printf("💾 Delegation saved to: %s\n", filename)
			return
		case "--raw":
			// Save raw CAR data for testing
			err := os.WriteFile("delegation.car", carData, 0644)
			if err != nil {
				log.Fatal("Failed to save CAR file:", err)
			}
			fmt.Printf("💾 Raw CAR saved to: delegation.car\n")
			return
		case "--debug":
			fmt.Printf("Issuer: %s\n", alice.DID().String())
			fmt.Printf("Audience: %s\n", bob.DID().String())
			fmt.Printf("Capabilities: %d\n", len(capabilities))
			fmt.Printf("CAR size: %d bytes\n", len(carData))
			fmt.Printf("Token size: %d chars\n", len(token))
			fmt.Printf("Token: %s\n", token)
			return
		}
	}

	// Show delegation info
	fmt.Printf("✅ UCAN Delegation Generated using go-ucanto!\n")
	fmt.Printf("═══════════════════════════════════════════\n")
	fmt.Printf("📊 Details:\n")
	fmt.Printf("   Issuer:       %s\n", alice.DID().String())
	fmt.Printf("   Audience:     %s\n", bob.DID().String())
	fmt.Printf("   Capabilities: %d\n", len(capabilities))
	
	// Handle expiration properly
	if exp := del.Expiration(); exp != nil {
		fmt.Printf("   Expires:      %s\n", time.Unix(int64(*exp), 0).Format(time.RFC3339))
	} else {
		fmt.Printf("   Expires:      Never\n")
	}
	
	fmt.Printf("   CAR Size:     %d bytes\n", len(carData))
	fmt.Printf("   Token Size:   %d chars\n", len(token))
	fmt.Printf("\n")

	// Show capabilities
	fmt.Printf("🔑 Capabilities:\n")
	for i, cap := range capabilities {
		fmt.Printf("   %d. %s on %s\n", i+1, cap.Can(), cap.With())
	}
	fmt.Printf("\n")

	// Show token (truncated if too long)
	fmt.Printf("🎯 Base64 CAR Token:\n")
	fmt.Printf("─────────────────────────────────────────\n")
	if len(token) > 200 {
		fmt.Printf("%s...\n", token[:200])
		fmt.Printf("(truncated - %d total chars)\n", len(token))
	} else {
		fmt.Printf("%s\n", token)
	}
	fmt.Printf("─────────────────────────────────────────\n")
	fmt.Printf("\n")

	// Show test commands
	fmt.Printf("📋 Test Commands:\n")
	fmt.Printf("\n")
	fmt.Printf("# Set token as environment variable:\n")
	fmt.Printf("export UCAN_TOKEN='%s'\n", token)
	fmt.Printf("\n")
	fmt.Printf("# Test health check:\n")
	fmt.Printf("curl -X GET http://localhost:8080/health | jq\n")
	fmt.Printf("\n")
	fmt.Printf("# Test parse delegation:\n")
	fmt.Printf("curl -X POST http://localhost:8080/api/parse/delegation \\\n")
	fmt.Printf("  -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("  -d \"$(jq -n --arg token \\\"$UCAN_TOKEN\\\" '{\\\"token\\\": $token}')\" | jq\n")
	fmt.Printf("\n")
	fmt.Printf("# Test parse chain:\n")
	fmt.Printf("curl -X POST http://localhost:8080/api/parse/chain \\\n")
	fmt.Printf("  -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("  -d \"$(jq -n --arg token \\\"$UCAN_TOKEN\\\" '{\\\"token\\\": $token}')\" | jq\n")
	fmt.Printf("\n")
	fmt.Printf("# Test validation:\n")
	fmt.Printf("curl -X POST http://localhost:8080/api/validate/chain \\\n")
	fmt.Printf("  -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("  -d \"$(jq -n --arg token \\\"$UCAN_TOKEN\\\" '{\\\"token\\\": $token}')\" | jq\n")
	fmt.Printf("\n")
	fmt.Printf("# Test graph generation:\n")
	fmt.Printf("curl -X POST http://localhost:8080/api/graph/delegation \\\n")
	fmt.Printf("  -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("  -d \"$(jq -n --arg token \\\"$UCAN_TOKEN\\\" '{\\\"token\\\": $token}')\" | jq\n")
	fmt.Printf("\n")
	fmt.Printf("# Test invocation graph:\n")
	fmt.Printf("curl -X POST http://localhost:8080/api/graph/invocation \\\n")
	fmt.Printf("  -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("  -d \"$(jq -n --arg token \\\"$UCAN_TOKEN\\\" '{\\\"token\\\": $token}')\" | jq\n")
	fmt.Printf("\n")
	fmt.Printf("💡 Options:\n")
	fmt.Printf("  --token-only     Output only the token\n")
	fmt.Printf("  --save [file]    Save delegation to file\n")
	fmt.Printf("  --raw            Save raw CAR file\n")
	fmt.Printf("  --debug          Show debug information\n")
	fmt.Printf("\n")
}