# Gift Grab Server

The [Nakama](https://heroiclabs.com/nakama/) game server for Gift Grab. It's a Go runtime module that signs players in with Clerk, adds a few custom RPCs for accounts, friends and groups, and sets up a monthly leaderboard. It runs in Docker next to Postgres, with Caddy in front for HTTPS on `gift-grab-server.app`.

## What the module adds

`main.go` registers everything in `InitModule`:

### Clerk sign-in (`controllers/auth/auth.go`)

A `BeforeAuthenticateCustom` hook. The client sends a Clerk session JWT as the custom auth ID. The hook verifies the RS256 signature against Clerk's JWKS (cached for an hour), then replaces the account ID and username with the Clerk user's ID and username. The JWKS URL depends on the `ENVIRONMENT` runtime variable:

| `ENVIRONMENT` | Clerk instance |
| --- | --- |
| `development` | the Clerk development instance (`*.clerk.accounts.dev`) |
| `production` | `clerk.gift-grab-server.app` |

Any other value fails sign-in.

### RPCs

| RPC | File | What it does |
| --- | --- | --- |
| `account_delete_id` | `controllers/account/account.go` | deletes the calling user's account |
| `get_friendship_state` | `controllers/friend/friend.go` | returns the friendship state between the caller and another user (`-1` if none) |
| `get_group_by_id` | `controllers/group/group.go` | returns a group as clean JSON |
| `get_group_membership_state` | `controllers/group/group.go` | returns the caller's membership state in a group |

Payloads are JSON, parsed and returned through `utils/utils.go`.

### Leaderboard (`controllers/leaderboard/leaderboard.go`)

`monthly_leaderboard`: descending, keeps each player's best score, resets at midnight UTC on the 1st of each month.

## Configuration

`local.yml` is copied into the image as Nakama's config. It sets debug logging, session lengths (1 hour access, 24 hour refresh), the Apple sign-in bundle ID, the console login and `ENVIRONMENT`. Replace `<CONSOLE_USERNAME>` and `<CONSOLE_PASSWORD>` before deploying.

## Running locally

Requires Docker with Compose.

```sh
docker compose up --build
```

This builds the module as a Go plugin (`Dockerfile`, Nakama plugin builder 3.16.0, vendored dependencies), runs Postgres 12.2, runs Nakama's migrations and starts the server.

| Port | Service |
| --- | --- |
| 7349 | gRPC API |
| 7350 | HTTP API |
| 7351 | console |
| 5432 | Postgres |

## Deploying

- `scripts/deploy-docker-hub.sh` asks for an environment (production, staging or development), writes it into `local.yml`, builds a `linux/arm64` image and pushes it to Docker Hub as `treycodes/gift-grab:latest`. The auth hook only knows `development` and `production`, so a `staging` build will reject every sign-in until a staging Clerk URL is added.
- `Caddyfile` serves `gift-grab-server.app`, sending `/v2/console*` and the console UI to port 7351 and the game API (`/v2*`) to port 7350.
- `docker-cleanup.sh` removes old containers and images on the host.

## Tech stack

- Go 1.20, `nakama-common` v1.27, `golang-jwt/jwt` v5
- Nakama 3.16, Postgres 12.2, Docker Compose, Caddy
- Clerk for player identity

## Related

- [gift_grab_client](https://github.com/trey-a-hope/gift_grab_client), [gift_grab_game](https://github.com/trey-a-hope/gift_grab_game), [gift_grab_ui](https://github.com/trey-a-hope/gift_grab_ui)
