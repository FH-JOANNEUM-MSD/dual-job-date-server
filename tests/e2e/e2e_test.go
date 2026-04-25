package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHealthEndpoint_E2E(t *testing.T) {
	baseURL := strings.TrimSpace(os.Getenv("E2E_BASE_URL"))
	if baseURL == "" {
		t.Skip("set E2E_BASE_URL to run e2e tests, e.g. http://localhost:8080")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(strings.TrimRight(baseURL, "/") + "/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed reading response body: %v", err)
	}

	if !strings.Contains(string(body), "Server läuft!") {
		t.Fatalf("unexpected health response: %q", string(body))
	}
}
