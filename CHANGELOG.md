# Changelog — grounded

## v0.1.0 — unreleased

- `groundedd`: skeleton for Grounded (`01006F301AE9C000`, UE4 project Maine).
  Env-driven config, `GET /healthz`, and stub `/v1/xbox/*` + `/v1/party/relay`
  endpoints answering 501 until the XAL/XSTS emulation pieces land. No NEX
  gameplay server applies (linked core only); relay/voice defers to the playfab
  Party relay once this title's TitleId is known.
