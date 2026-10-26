# Grounded

OpenPak's Grounded Server (stub).

| | |
| --- | --- |
| Game | Grounded |
| Network | Nintendo Switch |
| Title ID | `01006F301AE9C000` |
| Engine | Unreal Engine 4 (4.26+, project `Maine`) |
| Licence | AGPL-3.0-only |

## What it serves

Grounded on Switch plays online through Xbox Live auth plus a PlayFab Party
relay — there is no NEX gameplay server (see [Protocol](docs/protocol.md)):

| Host | What | Served by |
| --- | --- | --- |
| `*.xboxlive.com` | XAL/XSTS auth, matchmaking, session directory, invites, presence | this repo (`groundedd`) — stub, returns 501 until the pieces in [Xbox gate](docs/xbox-gate.md) exist |
| PlayFab Party 1.8.5 relay | DTLS relay + voice/chat transport (`Party.nro`) | [`playfab`](https://github.com/openpak/playfab) relay, once it serves this title |

## Run

```sh
go run ./cmd/groundedd
```

| Variable | Meaning |
| --- | --- |
| `GROUNDED_ADDR` | stub API listen address (default `:23560`) |
| `GROUNDED_HEALTH_ADDR` | `GET /healthz` listen address (default `:23561`) |
| `GROUNDED_PARTY_RELAY_URL` | PlayFab Party relay upstream (unset: relay stubbed) |
| `GROUNDED_PLAYFAB_TITLE` | PlayFab TitleId (outstanding — see docs) |
| `GROUNDED_SCID` | Xbox Service Configuration ID (outstanding) |
| `NX_INTERNAL_URL` | nx-baas internal API (default `http://127.0.0.1:20070`) |
| `NX_INTERNAL_KEY` | its `X-Internal-Key` (`OPENPAK_INTERNAL_KEY` is a fallback) |

## Docs

- [Protocol](docs/protocol.md)
- [Xbox gate](docs/xbox-gate.md)
- [Ports](docs/ports.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
