# zauth

Lightweight authentication server, with basic LDAP support, and a web UI.

## Compile

The build process requires:

- Git
- Go 1.13+
- MySQL 5.7+
- Task
- Goimports
- Minify

We will leave installing Go and MySQL to you. Installing Task, Goimports, and
Minify can be done like this:

```sh
go get -u github.com/go-task/task/v2/cmd/task
go get -u golang.org/x/tools/cmd/goimports
go get -u github.com/tdewolff/minify/cmd/minify
```

After that you can now build and run the application by using Task from the
command line

```sh
task build
./zauth
```

Or you can have Task build and run the app for you. This is nice, because it
will recompile and run the new application if you change any source files:

```sh
task run
```

## Docker Development

The repo includes a `docker-compose.yml` with two services: `mariadb` (the
database, with the schema loaded on first start) and `dev` (a Go toolchain
container with the repo mounted at `/src`).

```sh
# Start MariaDB and the dev container
docker compose up -d

# Seed the database with dev users and groups (fresh DB only)
docker compose exec dev go run ./cmd/seed

# Run the tests
docker compose exec dev go test ./...

# Build the binary (lands in the repo on the host via the bind mount)
docker compose exec dev go build -o zauth ./cmd
```

Database settings for tests and the seed command come from the `ZAUTH_DB_*`
environment variables (see `db.ConfigFromEnv`). The compose file sets them for
the `dev` container; on the host they default to `127.0.0.1:3306`, which the
`mariadb` service publishes, so a host-built binary works too.

To wipe the database and start over (re-runs the schema and lets you re-seed):

```sh
docker compose down -v && docker compose up -d
```

## FAQ

### I get "cross-origin request detected" (403) when submitting a form!

CSRF protection compares the browser's `Sec-Fetch-Site` and `Origin` headers
against the request's `Host`. A 403 means the browser considered the request
cross-origin, which usually means you reached zauth through a different
hostname or port than the page was served from. If zauth sits behind a reverse
proxy, make sure the proxy forwards the original `Host` header rather than
rewriting it to the backend address.

This is no longer affected by the `production` setting in your config, and no
longer depends on SSL/TLS.

### How can I query and test the LDAP server?

One way is to install `ldapsearch` which is standards compliant. Then you can:

```sh
# Get info about the user joshz (anonymously)
$ ldapsearch -x -H ldap://localhost:3389 "uid=joshz"
# Get info about the group admin (anonymously)
$ ldapsearch -x -H ldap://localhost:3389 "cn=admin"
# ldapwhoami currently doesn't work with zauth (BUG), so you can use to login
$ ldapsearch -x -H ldap://localhost:3389 -W -D 'uid=joshz'
#
# Login as joshz (and return all of your info) (currently broken)
$ ldapwhoami -x -H ldap://localhost:3389 -W -D 'uid=joshz'

```
