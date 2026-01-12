package main

import (
	"encoding/base64"
	"encoding/json"

	"github.com/skip2/go-qrcode"
)

func encodeSignedContentAsQRCode(signedContent SignedContent) (string, error) {
	jsonBytes, err := json.Marshal(signedContent)
	if err != nil {
		return "", err
	}

	return encodeStringAsQRCode(string(jsonBytes))
}

func encodeCOSEAsQRCode(coseBase64 string) (string, error) {
	return encodeStringAsQRCode(coseBase64)
}

func encodeStringAsQRCode(content string) (string, error) {
	qrCode, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return "", err
	}

	pngBytes, err := qrCode.PNG(256)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(pngBytes)
	return base64String, nil
}

