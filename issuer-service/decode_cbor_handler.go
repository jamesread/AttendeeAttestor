package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fxamacker/cbor/v2"
	"github.com/sirupsen/logrus"
)

type DecodeCBORRequest struct {
	Data string `json:"data"`
}

type DecodeCBORResponse struct {
	Decoded interface{} `json:"decoded"`
	Error   string      `json:"error,omitempty"`
}

func createDecodeCBORHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method": request.Method,
			"path":   request.URL.Path,
		}).Info("Decode CBOR request received")

		if request.Method != http.MethodPost {
			logrus.Warn("Invalid method for decode-cbor endpoint")
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req DecodeCBORRequest
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			logrus.WithError(err).Error("Failed to decode request body")
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			response := DecodeCBORResponse{
				Error: fmt.Sprintf("Invalid request data: %v", err),
			}
			json.NewEncoder(writer).Encode(response)
			return
		}

		if req.Data == "" {
			logrus.Warn("Empty data in decode-cbor request")
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			response := DecodeCBORResponse{
				Error: "Data is required",
			}
			json.NewEncoder(writer).Encode(response)
			return
		}

		cborBytes, err := decodeBase64Input(req.Data)
		if err != nil {
			logrus.WithError(err).Error("Failed to decode base64 input")
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			response := DecodeCBORResponse{
				Error: fmt.Sprintf("Failed to decode base64 input: %v", err),
			}
			json.NewEncoder(writer).Encode(response)
			return
		}

		var decoded interface{}
		if err := cbor.Unmarshal(cborBytes, &decoded); err != nil {
			logrus.WithError(err).Error("Failed to unmarshal CBOR")
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			response := DecodeCBORResponse{
				Error: fmt.Sprintf("Failed to decode CBOR: %v", err),
			}
			json.NewEncoder(writer).Encode(response)
			return
		}

		jsonSerializable := convertToJSONSerializable(decoded)

		logrus.Info("CBOR decoded successfully")
		response := DecodeCBORResponse{
			Decoded: jsonSerializable,
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			logrus.WithError(err).Error("Failed to encode response")
		}
	})
}

func decodeBase64Input(input string) ([]byte, error) {
	input = trimWhitespace(input)

	if bytes, err := base64.RawURLEncoding.DecodeString(input); err == nil && len(bytes) > 0 {
		return bytes, nil
	}

	if bytes, err := base64.URLEncoding.DecodeString(input); err == nil && len(bytes) > 0 {
		return bytes, nil
	}

	if bytes, err := base64.StdEncoding.DecodeString(input); err == nil && len(bytes) > 0 {
		return bytes, nil
	}

	return nil, io.ErrUnexpectedEOF
}

func trimWhitespace(s string) string {
	result := ""
	for _, r := range s {
		if r != ' ' && r != '\n' && r != '\r' && r != '\t' {
			result += string(r)
		}
	}
	return result
}

func convertToJSONSerializable(data interface{}) interface{} {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			keyStr := fmt.Sprintf("%v", k)
			result[keyStr] = convertToJSONSerializable(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = convertToJSONSerializable(val)
		}
		return result
	case []byte:
		return base64.StdEncoding.EncodeToString(v)
	default:
		return v
	}
}

