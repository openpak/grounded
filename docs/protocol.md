# Grounded (Switch) — network protocol

Wire facts read out of the Switch build (`01006F301AE9C000`, UE4 project
`Maine`, SDK 16.2.6.0). Keep this file in step with `cmd/groundedd`.

Recon sources (private, not shipped here):

- `scratch/decompile-grounded/RECON.md`
- `scratch/decompile-grounded/findings.json`
- `scratch/decompile-grounded/strings.txt` (curated; raw dumps are
  `strings-main.txt`, `strings-party.txt`, `strings-subsdk0.txt`)

## 1. Middleware verdict

**Xbox Live (XSAPI + XAL) for auth/session/matchmaking, PlayFab Party 1.8.5
for relay/voice. No NEX gameplay server, no Photon client, no EOS.**

- Xbox: 194 `Xbl*` + 507 `Xal*` references — `XblContext`, `XblHttpCall`,
  `matchmaking_service`, `XblGameInvite`, `XblMultiplayerMatchmakingServer`,
  `XblSocialRelationship`, XSTS sandbox handling; `OnlineSubsystemXSAPI` +
  `OnlineSubsystemSwitch` both present (`RECON.md` §"Network surface", item 1).
- Party: 489 `Party*` references in main (`PartyEndpointSendMessage`,
  `PartyNetworkCreateEndpoint`, chat controls, `PartyDevice*`); RomFS ships
  `nro/Party.nro` (6.7 MB Party runtime: DTLS endpoints, direct/relay links,
  `…/voices/list` TTS endpoint); build tag
  `PlayFabParty.Cpp.Switch-16.2.3`; `OnlineSubsystemPlayFab` present
  (`RECON.md`, item 2).
- `ClientGetPlayFabIDsFromNintendoSwitchDeviceIds` ⇒ Switch-device-ID →
  PlayFab login (`RECON.md`, item 2).

## 2. Hosts the client dials

From decompressed strings (`RECON.md`, item 1):

- `auth` / `notify` / `social` / `sessiondirectory` / `userpresence` /
  `vortex-events` under `.xboxlive.com`, plus OneCollector/vortex/settings
  telemetry endpoints.
- Transport: libcurl 7.64.1 (NintendoSDK 16.2.4) + nghttp2; `nn::socket`
  (`bsd:u`) init/finalize strings; no bundled PEM/cert paths in strings
  (`RECON.md`, item 5).

npdm net services: `bsd:u bsdcfg nifm:u nsd:u sfdnsres ssl ldn:u friend:u
bcat:u mnpp:app ngc:u ngct:u`; no `pgl`, no `dauth`; `Allow Debug: YES`
(`RECON.md`, item 5; `findings.json` → "npdm net services").

## 3. What is NOT here

- **Photon: not used.** Zero `ExitGames|PhotonPeer|LoadBalancing` hits in
  main/Party/subsdk0; only `GetPhotonAuthenticationToken` PlayFab API
  operation names (`RECON.md`, item 3; `findings.json` → "Photon not used").
- **NEX: linked core only.** 10 `/ws/NEX4.6/Pack/7923402a/OnlineCore` source
  paths (Core/Buffer/Scheduler/Result + PRUDPEndPoint); zero NEX game-protocol
  strings (no DataStore/MatchMaking/Gathering/Ranking/Utility/Messaging) and
  zero `npln` hits in main — PRUDP transport pulled in via the nn SDK static
  lib. No NEX server applies and no `switch-nex` registry work is needed
  (`RECON.md`, item 4; `findings.json` → "NEX linked core only").
- **EOS: absent** — 0 hits (`RECON.md`, item 5).
- subsdk0 is Nintendo SDK multimedia (OMX.Nvidia/NvMM/Android media), not net
  middleware (`RECON.md`, item 5; `findings.json` → "subsdk0 identity").

## 4. Outstanding keys

Exact UE 4.x minor, PlayFab TitleId, SCID, and sandbox live in the pak configs
(`Maine/Content/Paks/pakchunk0-Switch.*`, 2.6 GB IO Store; paks may be
encrypted — needs UnrealPak/unpak) or a runtime XAL log capture
(`RECON.md` §"Follow-ups"; `findings.json` → unresolved). No runtime traffic
capture (Harbor) yet to confirm secondary socket use (`ldn:u` + `bsd:u` are
present). Until the TitleId/SCID land, `groundedd` answers its stub API with
501 (see [Xbox gate](xbox-gate.md)).
