# Stockyard Scout

**Dead link and SSL checker — crawl a domain, report broken links and expired certs on a schedule**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9100:9100 -v scout_data:/data ghcr.io/stockyard-dev/stockyard-scout
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9100` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9100` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `SCOUT_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 2 sites, weekly schedule | Unlimited sites, hourly schedule |
| Price | Free | $2.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Developer Tools

## License

Apache 2.0
