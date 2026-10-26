# Ports

Claimed from the game-services range (`23500–23999`, one block of ten per
title), after unrailed (23500–23509), fall-guys (23530–23539),
plants-vs-zombies (23540–23549) and stumble-guys (23550–23559).

| Port | Use |
| --- | --- |
| 23560 | plain HTTP: Xbox-gate stub API (behind the TLS terminator) |
| 23561 | `/healthz` |
| 23562–23569 | reserved for this title (XSTS/session-directory listeners, once built) |

Relay/voice transport rides the playfab service's ports, not this block.
