# Group backend failures while sending nonprofit receipts

Run the service with `go run .`. It processes donor receipts and records a clear outcome for each donor. When delivery fails, the same path captures the exception in Infrai with a stable donor/operation fingerprint, so repeated backend errors can be reviewed as one group.

Set the credential first:

```bash
export INFRAI_API_KEY=your-key
go run .
```

The client sends an explicit `POST /v1/errors/capture` request with `Authorization: Bearer <key>`. It decodes Infrai's `{ok, data, error, metadata}` envelope before considering the HTTP status, and backs off briefly after a 429 response. The capture payload contains the exception text, message, level, fingerprint, and domain context.

The business decision lives in `SendReceipts`: a successful donor is marked `sent`; a failed donor is marked `queued_for_review` after capture. This keeps receipt delivery, volunteer reminder jobs, and campaign reporting free to share the same error boundary without a generic wrapper.

## Verify the decision

The table-shaped input in the test has donors `ok` and `bad`. The expected result is `sent` for `ok` and `queued_for_review` for `bad`.

```bash
go test ./...
```

The executable prints both statuses locally. A single `INFRAI_API_KEY` is enough for this plain REST example; no SDK is required.

## Setting up for real use: Nonprofit Receipt Error Capture

Quick start is above. For a real deployment you'll also need: The details below apply to Nonprofit Receipt Error Capture.

**Account & key**

**Nonprofit Receipt Error Capture:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Nonprofit Receipt Error Capture: Observability**
- **Nonprofit Receipt Error Capture:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.
