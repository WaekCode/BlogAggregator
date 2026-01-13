# gator — Blog Aggregator CLI 🐊

**gator** (a.k.a. BlogAggregator) is a small CLI for registering users, following RSS feeds and browsing collected posts.

---

## 🔧 Prerequisites

- **PostgreSQL** (local or remote) — used to store users, feeds and posts.
- **Go** (1.25 or newer) — for building or installing the CLI.

---

## 📦 Install the CLI

Option A — install from the module with `go install` (recommended):

```bash
# installs the main package binary into $GOBIN or $GOPATH/bin
go install github.com/WaekCode/gator@latest
```

Make sure your Go bin directory is on your PATH (for example: `export PATH="$PATH:$(go env GOPATH)/bin"`).

Option B — build a `gator` binary locally and place it in your PATH:

```bash
# build a local binary named 'gator'
go build -o gator .
# move it to somewhere in your PATH
mv gator /usr/local/bin/
```

---

## ⚙️ Configuration

The CLI expects a JSON config file at `~/.gatorconfig.json`.

Create `~/.gatorconfig.json` with at least the `db_url` field. Example:

```json
{
  "db_url": "postgres://postgres:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

- `db_url` — PostgreSQL connection string used by the app.
- `current_user_name` — optional; can be set by `gator register`/`gator login`.

---

## 🗄️ Database setup

Create a PostgreSQL database and apply the schema files in `sql/schema/` in numeric order:

```bash
createdb gator
psql -d gator -f sql/schema/001_users.sql
psql -d gator -f sql/schema/002_feed.sql
psql -d gator -f sql/schema/003_feedFollowers.sql
psql -d gator -f sql/schema/004_last_fetched_at.sql
psql -d gator -f sql/schema/005_posts.sql
```

Adjust the DB name, user and password as needed and ensure the `db_url` in `~/.gatorconfig.json` matches.

---

## ▶️ Running the program

If you installed via `go install`, run the binary (binary name will follow the package name, or use the local `gator` binary from Option B):

```bash
# show available commands
gator commands
# or when running directly from source
go run main.go commands
```

### Common commands

- `gator register <username>` — create a new user and set them as current
- `gator login <username>` — set the current user
- `gator users` — list users (current user is marked)
- `gator addfeed "Name" "https://example.com/rss"` — add a feed (must be logged in)
- `gator feeds` — list all feeds
- `gator follow <feed_url>` — follow a feed (must be logged in)
- `gator following` — list feeds the current user follows
- `gator agg <duration>` — run the aggregator periodically (e.g., `gator agg 10s`)
- `gator browse [limit]` — browse collected posts (optional limit)

---

## 🛠️ Troubleshooting & tips

- If the CLI fails with "no such file or directory" for the config, create `~/.gatorconfig.json` as shown above.
- Double-check that Postgres is running and your `db_url` is correct.
- If you prefer, run `go run main.go <command>` while developing.

---