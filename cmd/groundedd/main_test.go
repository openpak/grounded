package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStubEndpointsAnswer501(t *testing.T) {
	s := &server{xbox: newXboxStub(Config{})}
	for _, target := range []string{"/v1/xbox/auth", "/v1/xbox/session", "/v1/party/relay"} {
		var h http.HandlerFunc
		switch target {
		case "/v1/xbox/auth":
			h = s.xboxAuth
		case "/v1/xbox/session":
			h = s.xboxSession
		default:
			h = s.partyRelay
		}
		w := httptest.NewRecorder()
		h(w, httptest.NewRequest(http.MethodGet, target, nil))
		if w.Code != http.StatusNotImplemented {
			t.Fatalf("%s: got %d, want 501", target, w.Code)
		}
		var body map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: response is not JSON: %v", target, err)
		}
		if body["error"] != "not_implemented" || body["missing"] == "" {
			t.Fatalf("%s: body %v names no missing piece", target, body)
		}
	}
}

func TestConfigDefaults(t *testing.T) {
	t.Setenv("GROUNDED_ADDR", "")
	t.Setenv("GROUNDED_HEALTH_ADDR", "")
	t.Setenv("NX_INTERNAL_URL", "")
	cfg := configFromEnv()
	if cfg.Addr != ":23560" || cfg.HealthAddr != ":23561" {
		t.Fatalf("default addrs: got %q %q", cfg.Addr, cfg.HealthAddr)
	}
	if cfg.NXInternalURL != "http://127.0.0.1:20070" {
		t.Fatalf("default nx-baas url: got %q", cfg.NXInternalURL)
	}
}
