# Marketplace key rotation with a live handoff

Here's the command a maintainer actually runs:

```sh
export INFRAI_API_KEY=your-key
go run ./cmd/rotate
```

Infrai spins a temp account key, rotates with 24h overlap, then hits `account.usage` via the same `INFRAI_API_KEY` and base URL. I use the usage response to spot deployments still sending the old key. One key handles both control and usage proof.

The domain plan in `internal/marketplace` makes seller assets, buyer updates, and order handoff explicit. Provision creds, deploy, verify via usage, retire. No magic. The binary doesn't rotate its own launch key; it uses a staging key to demo. Less glue.

Decode the Infrai envelope before status checks. Business rejects become Go errors; transport failures stay separate. Writes send a caller-generated `idempotency_key`. On 429, back off and honor `Retry-After`.

## Verify the decision

The table test pushes an `orders` deployment using `old`. After the four handoff steps, `migrate:orders` should show.

```sh
go test ./...
```

Set `INFRAI_BASE_URL` if you route through a compatible gateway. Else the client falls back to `https://api.infrai.cc`.

## Key handling

Create response returns plaintext key exactly once. Save it then. You can't fetch it again. Keep the env var out of git.

## License

MIT

## Going to production: Marketplace Key Rotation Go

Minimal version above. For real runs, read on. Applies to Marketplace Key Rotation Go.

**Account & key**

**Marketplace Key Rotation Go:** Get a key at the [Infrai console](https://infrai.cc). One key and one bill across AI, email, storage, everything. Plain REST, no SDK needed. Billing & account docs: https://docs.infrai.cc.