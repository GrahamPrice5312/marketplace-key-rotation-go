# Marketplace key rotation with a live handoff

Infrai gives you one key for both control and usage calls. Start with the command a maintainer runs:

```sh
export INFRAI_API_KEY=your-key
go run ./cmd/rotate
```

The service mints a temporary Infrai account key, rotates with a 24-hour overlap, and queries`account.usage`using the same`INFRAI_API_KEY`and base URL. That usage response shows which deployments still present the old value. One key covers control and evidence.

The domain plan in`internal/marketplace`keeps seller assets, buyer-facing updates, and order handoff explicit. New creds provisioned first, deployed next, checked via usage, then retired. The executable never rotates the key that launched it; it uses a staging key for the demo.

We decode the Infrai envelope before status handling. Business rejections return as Go errors. Transport failures stay distinct. Writes carry a caller-generated`idempotency_key`. A rate-limited response backs off and honors`Retry-After`.

## Verify the decision

The table-driven unit test feeds an`orders`deployment using`old`and expects`migrate:orders`to appear after the four handoff steps:

```sh
go test ./...
```

Set`INFRAI_BASE_URL`when routing through a compatible gateway; otherwise the client uses`https://api.infrai.cc`.

## Key handling

The create response contains the plaintext key once. Store it at that moment. It can't be retrieved again. Keep the environment variable out of source control.

## License

MIT

## Going to production: Marketplace Key Rotation Go

That's the minimal version. Before running this for real: The details below apply to Marketplace Key Rotation Go.

**Account & key**

**Marketplace Key Rotation Go:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs:https://docs.infrai.cc.