package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

type SignedContent struct {
	Data      map[string]interface{} `json:"data"`
	Signature string                 `json:"signature"`
}

func generatePrivateKey() (ed25519.PrivateKey, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func signJSONContent(data map[string]interface{}, privateKey ed25519.PrivateKey) (SignedContent, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return SignedContent{}, err
	}

	signature := ed25519.Sign(privateKey, jsonBytes)
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	return SignedContent{
		Data:      data,
		Signature: signatureBase64,
	}, nil
}

func verifySignature(data map[string]interface{}, signatureBase64 string, publicKey ed25519.PublicKey) bool {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return false
	}

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false
	}

	return ed25519.Verify(publicKey, jsonBytes, signature)
}

func loadOrGeneratePrivateKey() (ed25519.PrivateKey, error) {
	return generatePrivateKey()
}

