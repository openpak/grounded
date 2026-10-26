package main

// xbox.go — middleware client stub for Grounded's network stack.
//
// The Switch client reaches Xbox Live through XSAPI calls backed by XAL
// (auth/notify/social/sessiondirectory/userpresence/vortex-events .xboxlive.com
// hosts) and plays over a PlayFab Party 1.8.5 DTLS relay (Party.nro). Neither
// upstream exists in OpenPak yet, so this file holds the seam both will plug
// into: the stub records what each endpoint needs and answers 501 naming the
// missing piece. See docs/xbox-gate.md for the emulation work behind it and
// docs/protocol.md for the recon evidence.

import (
	"log"
	"net/http"
)

// xboxStub is the placeholder for the XAL/XSTS emulation client. It carries
// the configured upstreams so the handlers already read from one place.
type xboxStub struct {
	cfg Config
}

func newXboxStub(cfg Config) *xboxStub { return &xboxStub{cfg: cfg} }

// xboxAuth stands in for the XAL/XSTS token endpoint the client must clear
// before anything else: without an Xbox user token there is no matchmaking,
// session directory, invite, or presence. Requires GROUNDED_SCID plus the XSTS
// emulator (docs/xbox-gate.md §2).
func (s *server) xboxAuth(w http.ResponseWriter, r *http.Request) {
	log.Printf("[grounded] xboxAuth from %s: stub, no XSTS emulator", r.RemoteAddr)
	writeStub(w, "xsts-emulator (GROUNDED_SCID unset or emulator unbuilt)")
}

// xboxSession stands in for the Xbox matchmaking/session-directory calls
// (XblMultiplayerMatchmakingServer, sessiondirectory.xboxlive.com). Requires a
// working xboxAuth plus the session-directory emulator (docs/xbox-gate.md §3).
func (s *server) xboxSession(w http.ResponseWriter, r *http.Request) {
	log.Printf("[grounded] xboxSession from %s: stub, no session-directory emulator", r.RemoteAddr)
	writeStub(w, "session-directory-emulator")
}

// partyRelay is the handoff to the PlayFab Party-compatible relay/voice
// transport (platforms/switch/games/playfab, internal/party). Requires
// GROUNDED_PARTY_RELAY_URL once that service serves this title's
// GROUNDED_PLAYFAB_TITLE; the TitleId is still outstanding.
func (s *server) partyRelay(w http.ResponseWriter, r *http.Request) {
	log.Printf("[grounded] partyRelay from %s: stub, relay=%q title=%q",
		r.RemoteAddr, s.xbox.cfg.PartyRelayURL, s.xbox.cfg.PlayFabTitleID)
	writeStub(w, "party-relay-upstream (GROUNDED_PARTY_RELAY_URL unset or title unconfigured)")
}
