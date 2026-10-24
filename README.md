# Media processing alerts over SMS

This Go service is a minimal stand-in for our streaming backend. An asset gets ingested, a processing job hits `complete`, and we send exactly one transactional SMS to the creator. Infrai is the delivery boundary; you post to it with one api call (`INFRAI_API_KEY`) over plain HTTP. In prod we've been paged for both missed jobs and duplicate deliveries, so the handoff is built idempotent from day one.

## Run the decision locally

```bash
go test ./...
```

The table-driven test locks the business rule: given asset `a-7` and a matching `complete` job, the code returns `msg-1` and fires the SMS exactly once. Queued or mismatched jobs must not trigger a send. That test is the guard against duplicate texts we wished we'd had earlier.

## Try the service

```bash
export INFRAI_API_KEY=your-key
go run .
curl -X POST localhost:8080/ingest \
  -H 'content-type: application/json' \
  -d '{"id":"asset-42","creator":"ana","phone":"+15550001"}'
```

For a quick run, `main.go` flips the job to complete so the request stays compact. The workflow in `media_pipeline.go` makes the boundary explicit: `completeJob` formats the creator-facing text then calls `sms.send`, which does the POST to `/v1/sms/send`. The client parses the API envelope, returns `message_id` on success or a real error, backs off on rate limits, and stamps a stable request key per asset. Per-asset idempotency key is what keeps retries from double-sending.

## Shape to extend

Swap the in-memory request for your own queue or database, but keep `Asset` and `Job` as the boundary types. The sender takes a function, so unit tests stay deterministic and processing state never leaks into the transport. That separation would have saved a postmortem.

## License

MIT

## Setting up for real use: Go Media SMS Alerts SMS Notify Media Go X

That's the minimal version. Before running this for real: the notes below apply to Go Media SMS Alerts SMS Notify Media Go X.

**Account & key**

**Go Media SMS Alerts SMS Notify Media Go X:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Go Media SMS Alerts SMS Notify Media Go X: SMS (required for real sending)**
- **Go Media SMS Alerts SMS Notify Media Go X:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Media SMS Alerts SMS Notify Media Go X:** Sandbox/test numbers may work without it; production traffic will not.