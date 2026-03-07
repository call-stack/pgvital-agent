# PGVitals Agent

A lightweight, open-source agent that collects PostgreSQL performance metrics and sends them to [PGVitals](https://pgvitals.kafal.studio) for analysis, alerting, and weekly digest reports.

## What it does

The agent connects to your PostgreSQL database as a **read-only** user and collects:

- **Query statistics** from `pg_stat_statements` (top 200 by total execution time)
- **Table statistics** from `pg_stat_user_tables` (scan counts, tuple counts, vacuum info)
- **Index statistics** from `pg_stat_user_indexes` (usage, sizes)

It then sends this data to the PGVitals API at a configurable interval (default: every 5 minutes).

## What it does NOT do

- Does not write to your database
- Does not store credentials on disk (only the instance registration ID)
- Does not collect actual row data or query parameters
- Does not require superuser access
- Query literals are scrubbed before transmission (`pg_stat_statements` normalizes them to `$1`, `$2`, etc., and the agent applies a second pass replacing any remaining literals with `?`)

## Quick start

### Docker (recommended)

```bash
docker run -d \
  -e PGVITALS_DATABASE_URL="postgres://pgvitals:your_password@host:5432/dbname" \
  -e PGVITALS_API_KEY="pgp_live_your_key_here" \
  --name pgvitals-agent \
  kalpitpant/pgvital-agent:latest
```

### From source

```bash
git clone https://github.com/pgvitals/agent.git
cd agent
go build -o pgvitals-agent ./cmd/pgvitals-agent

export PGVITALS_DATABASE_URL="postgres://pgvitals:your_password@host:5432/dbname"
export PGVITALS_API_KEY="pgp_live_your_key_here"

./pgvitals-agent
```

## Prerequisites

1. **PostgreSQL 12+** with `pg_stat_statements` enabled:

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
```

2. **A read-only database user** (recommended):

```sql
CREATE USER pgvitals WITH PASSWORD 'a_strong_password_here';
GRANT CONNECT ON DATABASE your_database TO pgvitals;
GRANT SELECT ON pg_stat_statements TO pgvitals;
GRANT USAGE ON SCHEMA public TO pgvitals;
```

The agent only needs SELECT access to system statistics views. Using a dedicated read-only user ensures it cannot modify any data even if compromised.

## Configuration

All configuration is done via environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PGVITALS_DATABASE_URL` | Yes | - | PostgreSQL connection string |
| `PGVITALS_API_KEY` | Yes | - | Your organization API key from the dashboard |
| `PGVITALS_SERVER_URL` | No | `https://api.pgvitals.kafal.studio` | PGVitals API server URL |
| `PGVITALS_INTERVAL` | No | `5m` | Collection interval (min: `1m`, max: `1h`) |
| `PGVITALS_HOSTNAME` | No | System hostname | Human-readable label for this instance |
| `PGVITALS_DATA_DIR` | No | `~/.pgvitals` | Directory to store instance registration |

## Architecture

```
pgvitals-agent
├── cmd/pgvitals-agent/     # Entry point
└── internal/
    ├── collector/           # SQL queries for pg_stat_statements, tables, indexes
    ├── config/              # Environment variable parsing
    ├── scrubber/            # Second-pass literal scrubbing
    └── sender/              # HTTPS transport to PGVitals API
```

## Data flow

```
PostgreSQL
  └─ pg_stat_statements ──┐
  └─ pg_stat_user_tables ─┤── Agent collects every 5m
  └─ pg_stat_user_indexes ┘       │
                                  ▼
                          Scrub any remaining
                          literal values
                                  │
                                  ▼
                          Send via HTTPS to
                          PGVitals API
                                  │
                                  ▼
                          Dashboard, alerts,
                          weekly digest email
```

## Security

- **Read-only**: The agent executes only SELECT queries against statistics views
- **No data exposure**: Query parameters are normalized by PostgreSQL and double-scrubbed by the agent
- **Minimal permissions**: Works with a restricted database user that has no access to your actual tables
- **Open source**: Inspect every SQL query the agent runs in [`internal/collector/`](internal/collector/)
- **Encrypted transport**: All data is sent over HTTPS

## Building the Docker image

```bash
docker build -t pgvitals-agent .

# Multi-platform build
docker buildx build --platform linux/amd64,linux/arm64 -t pgvitals-agent .
```

## License

MIT - see [LICENSE](LICENSE) for details.
