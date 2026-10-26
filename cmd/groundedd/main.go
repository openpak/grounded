// Command groundedd is OpenPak's Grounded (Switch, 01006F301AE9C000) server stub.
//
// Grounded's online path is Xbox Live (XSAPI+XAL auth, matchmaking, session
// directory) over a PlayFab Party 1.8.5 DTLS relay; there is no NEX gameplay
// server to implement (see docs/protocol.md). This skeleton owns the service's
// listeners and configuration and answers /healthz; the Xbox-gate endpoints and
// the Party-relay handoff are stubs that return 501 until the XAL/XSTS
// emulation pieces in docs/xbox-gate.md are built.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Config is the service's env-driven configuration.
type Config struct {
	// Addr serves the Xbox-gate stub API (behind the TLS terminator).
	Addr string
	// HealthAddr serves GET /healthz.
	HealthAddr string
	// PartyRelayURL is the PlayFab Party-compatible relay/voice upstream
	// (platforms/switch/games/playfab), once it serves this title.
	PartyRelayURL string
	// PlayFabTitleID and SCID come from the game's pak configs; both are still
	// outstanding (see docs/protocol.md). Empty means unconfigured.
	PlayFabTitleID string
	SCID           string
	// NXInternalURL/Key point at nx-baas for Nintendo id-token checks.
	NXInternalURL string
	NXInternalKey string
}

func configFromEnv() Config {
	return Config{
		Addr:           envOr("GROUNDED_ADDR", ":23560"),
		HealthAddr:     envOr("GROUNDED_HEALTH_ADDR", ":23561"),
		PartyRelayURL:  os.Getenv("GROUNDED_PARTY_RELAY_URL"),
		PlayFabTitleID: os.Getenv("GROUNDED_PLAYFAB_TITLE"),
		SCID:           os.Getenv("GROUNDED_SCID"),
		NXInternalURL:  envOr("NX_INTERNAL_URL", "http://127.0.0.1:20070"),
		NXInternalKey:  firstNonEmpty(os.Getenv("NX_INTERNAL_KEY"), os.Getenv("OPENPAK_INTERNAL_KEY")),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}

func main() {
	cfg := configFromEnv()
	s := &server{xbox: newXboxStub(cfg)}

	go func() {
		h := http.NewServeMux()
		h.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "groundedd"})
		})
		log.Printf("[grounded] health on %s", cfg.HealthAddr)
		log.Fatal(http.ListenAndServe(cfg.HealthAddr, h))
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/xbox/auth", s.xboxAuth)
	mux.HandleFunc("/v1/xbox/session", s.xboxSession)
	mux.HandleFunc("/v1/party/relay", s.partyRelay)
	log.Printf("[grounded] serving Xbox-gate stub API on %s (title=%q scid=%q relay=%q)",
		cfg.Addr, cfg.PlayFabTitleID, cfg.SCID, cfg.PartyRelayURL)
	log.Fatal(http.ListenAndServe(cfg.Addr, mux))
}

type server struct {
	xbox *xboxStub
}

func writeStub(w http.ResponseWriter, missing string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":   "not_implemented",
		"missing": missing,
	})
}
