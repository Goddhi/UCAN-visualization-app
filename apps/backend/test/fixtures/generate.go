package fixtures

import (
	"time"
	"io"

	"github.com/storacha/go-ucanto/core/delegation"
	"github.com/storacha/go-ucanto/principal/ed25519/signer"
	"github.com/storacha/go-ucanto/ucan"
)

// GenerateValidUCAN creates a valid UCAN delegation for testing
func GenerateValidUCAN() ([]byte, error) {
	// Create issuer (Alice)
	alice, err := signer.Generate()
	if err != nil {
		return nil, err
	}

	// Create audience (Bob)
	bob, err := signer.Generate()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// Archive as CAR bytes
	archive := del.Archive()
	return io.ReadAll(archive)
}

// GenerateExpiredUCAN creates an expired UCAN for testing
func GenerateExpiredUCAN() ([]byte, error) {
	alice, err := signer.Generate()
	if err != nil {
		return nil, err
	}

	bob, err := signer.Generate()
	if err != nil {
		return nil, err
	}

	// Create delegation that expired yesterday
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
		delegation.WithExpiration(int(time.Now().Add(-24*time.Hour).Unix())),
		delegation.WithNotBefore(int(time.Now().Add(-48*time.Hour).Unix())),
	)
	if err != nil {
		return nil, err
	}

	archive := del.Archive()
	return io.ReadAll(archive)
}

// GenerateComplexChain creates a multi-level delegation chain
func GenerateComplexChain() ([]byte, error) {
	// Root (Alice)
	alice, err := signer.Generate()
	if err != nil {
		return nil, err
	}

	// Bob
	bob, err := signer.Generate()
	if err != nil {
		return nil, err
	}

	// Charlie
	charlie, err := signer.Generate()
	if err != nil {
		return nil, err
	}

	// 1. Alice delegates to Bob
	aliceToBob, err := delegation.Delegate(
		alice,
		bob,
		[]ucan.Capability[ucan.NoCaveats]{
			ucan.NewCapability(
				"store/*",
				"storage:*",
				ucan.NoCaveats{},
			),
		},
		delegation.WithExpiration(int(time.Now().Add(7*24*time.Hour).Unix())),
	)
	if err != nil {
		return nil, err
	}

	// 2. Bob delegates to Charlie (with Alice->Bob as proof)
	bobToCharlie, err := delegation.Delegate(
		bob,
		charlie,
		[]ucan.Capability[ucan.NoCaveats]{
			ucan.NewCapability(
				"store/add",
				"storage:alice/*",
				ucan.NoCaveats{},
			),
		},
		delegation.WithExpiration(int(time.Now().Add(24*time.Hour).Unix())),
		delegation.WithProof(delegation.FromDelegation(aliceToBob)),
	)
	if err != nil {
		return nil, err
	}

	archive := bobToCharlie.Archive()
	return io.ReadAll(archive)
}