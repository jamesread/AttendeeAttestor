package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSHandlerSetsAllowOriginHeader(t *testing.T) {
	handler := createCORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodPost, "/generate", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	allowOrigin := recorder.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "http://localhost:3000" {
		t.Errorf("Expected Access-Control-Allow-Origin to be http://localhost:3000, got %s", allowOrigin)
	}
}

func TestCORSHandlerSetsAllowMethodsHeader(t *testing.T) {
	handler := createCORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodPost, "/generate", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	allowMethods := recorder.Header().Get("Access-Control-Allow-Methods")
	if allowMethods != "POST, OPTIONS" {
		t.Errorf("Expected Access-Control-Allow-Methods to be POST, OPTIONS, got %s", allowMethods)
	}
}

func TestCORSHandlerSetsAllowHeadersHeader(t *testing.T) {
	handler := createCORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodPost, "/generate", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	allowHeaders := recorder.Header().Get("Access-Control-Allow-Headers")
	if allowHeaders != "Content-Type" {
		t.Errorf("Expected Access-Control-Allow-Headers to be Content-Type, got %s", allowHeaders)
	}
}

func TestCORSHandlerHandlesOPTIONSRequest(t *testing.T) {
	handler := createCORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodOptions, "/generate", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS request, got %d", recorder.Code)
	}

	allowOrigin := recorder.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "http://localhost:3000" {
		t.Errorf("Expected Access-Control-Allow-Origin header in OPTIONS response")
	}
}

func TestCORSHandlerCallsNextHandlerForNonOPTIONSRequest(t *testing.T) {
	wasCalled := false
	handler := createCORSHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wasCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodPost, "/generate", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if !wasCalled {
		t.Error("Expected next handler to be called for POST request")
	}
}

