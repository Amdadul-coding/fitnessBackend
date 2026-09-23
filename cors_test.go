package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWorkoutBrowserPreflight(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/workouts", nil)
	request.Header.Set("Origin", "https://example.onrender.com")
	request.Header.Set("Access-Control-Request-Method", "POST")
	request.Header.Set("Access-Control-Request-Headers", "content-type")
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, request)
	if response.Code != http.StatusNoContent || response.Header().Get("Access-Control-Allow-Origin") != "*" || response.Header().Get("Access-Control-Allow-Methods") != "POST, OPTIONS" || response.Header().Get("Access-Control-Allow-Headers") != "Content-Type" {
		t.Fatalf("preflight failed: %d %v", response.Code, response.Header())
	}
}

func TestWorkoutBrowserPost(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(`{"bodyParts":"chest","workoutType":"home"}`))
	request.Header.Set("Origin", "https://example.onrender.com")
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("browser POST failed: %d %v", response.Code, response.Header())
	}
}
