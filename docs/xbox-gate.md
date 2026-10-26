# Xbox gate — auth flow and emulation pieces

Grounded's online path is gated on Xbox Live: XAL/XSTS auth plus Xbox
matchmaking/session directory are required, with no anonymous or
Nintendo-only fallback visible in the binary (`protocol.md` §1). This file
describes the flow the client expects and the server pieces needed to answer
it. Scope is comparable to the minecraft-dungeons Xbox gate, not a NEX title.

## 1. Auth flow (as the client drives it)

1. **XAL sign-in.** The client authenticates the player through XAL against
   `auth.xboxlive.com` and obtains Xbox user tokens (XSTS sandbox handling
   strings are present in main).
2. **Service ticket.** XSTS issues a token for the title's Service
   Configuration (SCID — outstanding, see `protocol.md` §4).
3. **PlayFab login.** `ClientGetPlayFabIDsFromNintendoSwitchDeviceIds` maps the
   Switch device id to a PlayFab identity (PlayFab TitleId outstanding).
4. **Session/matchmaking.** Xbox matchmaking (`XblMultiplayerMatchmakingServer`)
   and the session directory (`sessiondirectory.xboxlive.com`) place the player;
   invites (`XblGameInvite`), social graph (`XblSocialRelationship`), presence
   (`userpresence.xboxlive.com`), and notify/vortex-events channels ride the
   same Xbox user token.
5. **Relay/voice.** Gameplay traffic goes over the PlayFab Party 1.8.5 DTLS
   relay (`Party.nro`: endpoints, direct/relay links); see §4.

`groundedd` currently stubs steps 1–4 at `/v1/xbox/auth` and `/v1/xbox/session`
(HTTP 501 naming the missing piece); step 5 is handed to the Party relay
upstream at `/v1/party/relay`.

## 2. XAL/XSTS emulation pieces needed

- **XSTS emulator** — answers the `auth.xboxlive.com` token exchange with
  OpenPak-issued Xbox user tokens for the title's SCID (needs `GROUNDED_SCID`
  from pak configs or a runtime XAL log capture). Fronts `POST /v1/xbox/auth`.
- **Session-directory emulator** — minimal
  `sessiondirectory.xboxlive.com` + `XblMultiplayerMatchmakingServer`
  semantics: create/join/find sessions so hosts and clients rendezvous.
  Fronts `POST /v1/xbox/session`.
- **Social/presence/invite stubs** — `social`, `userpresence`, `notify`,
  `vortex-events` endpoints returning empty-but-valid shapes (friend list,
  presence, invite delivery) so the client never blocks on them.
- **Telemetry sink** — OneCollector/vortex/settings endpoints accepting and
  discarding events.
- **TLS story** — the client dials `https://*.xboxlive.com`; requests arrive at
  the gate behind whatever terminates that TLS (same pattern as the other game
  services: plain HTTP listener + terminator in front).

## 3. What is explicitly out of scope

- NEX gameplay protocols — linked core only, nothing to implement
  (`protocol.md` §3).
- Photon/EOS — absent from the binary (`protocol.md` §3).
- Game simulation — sessions are relayed, not simulated.

## 4. Party-relay reuse

Relay/voice transport should reuse `platforms/switch/games/playfab`
(`internal/party`: Party QoS beacon, relay handling) rather than reimplementing
Party DTLS here — `groundedd` only holds the handoff (`GROUNDED_PARTY_RELAY_URL`)
until that service serves this title's PlayFab TitleId. Blocked on the same
outstanding TitleId (`protocol.md` §4).
