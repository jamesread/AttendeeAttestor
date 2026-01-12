package main

import (
	"crypto/ed25519"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

const (
	COSEAlgorithmEd25519 = -8
)

type COSESign1 struct {
	Protected   []byte
	Unprotected map[int]interface{}
	Payload     []byte
	Signature   []byte
}

func createCOSESign1(payload map[string]interface{}, privateKey ed25519.PrivateKey, keyID string) ([]byte, error) {
	payloadCBOR, err := encodeToCBOR(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}

	protectedHeaders := map[int]interface{}{
		1: COSEAlgorithmEd25519,
		4: keyID,
	}

	protectedCBOR, err := encodeToCBOR(protectedHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to encode protected headers: %w", err)
	}

	unprotectedHeaders := map[int]interface{}{}

	sigStructure := []interface{}{
		"Signature1",
		protectedCBOR,
		[]byte{},
		payloadCBOR,
	}

	sigStructureCBOR, err := encodeToCBOR(sigStructure)
	if err != nil {
		return nil, fmt.Errorf("failed to encode signature structure: %w", err)
	}

	signature := ed25519.Sign(privateKey, sigStructureCBOR)

	coseSign1 := COSESign1{
		Protected:   protectedCBOR,
		Unprotected: unprotectedHeaders,
		Payload:     payloadCBOR,
		Signature:   signature,
	}

	coseSign1Array := []interface{}{
		coseSign1.Protected,
		coseSign1.Unprotected,
		coseSign1.Payload,
		coseSign1.Signature,
	}

	return encodeToCBOR(coseSign1Array)
}

func encodeToCBOR(data interface{}) ([]byte, error) {
	em, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return nil, err
	}

	return em.Marshal(data)
}

func signCOSE(payload map[string]interface{}, privateKey ed25519.PrivateKey, keyID string) ([]byte, error) {
	return createCOSESign1(payload, privateKey, keyID)
}

