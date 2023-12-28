# Cubari metrics

At its core, this is a Prometheus remote writer that ingests datapoints for counters on Cubari properties: `cubari.moe` and `proxy.cubari.moe`.

This telemetry is privacy-preserving (no fingerprinting) and is solely used to gauge proxy popularity without digging into GA4 (which can be blocked, and is hard to get aggregate data) or server-side instrumentation (which is inherently less accurate due to downstream caching).

## How to use

This is deployed with dokku, but the configs required are:

```
REMOTE_WRITE_URL
REMOTE_WRITE_USERNAME
REMOTE_WRITE_PASSWORD
PORT
```