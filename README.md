# leafnet

Web application for building and exploring genealogy trees, written in Go. People and their parent relationships are stored as a graph in Neo4j and rendered as an interactive 3D force-directed graph in the browser.

## Features

- User registration and authentication (bcrypt-hashed passwords, cookie sessions)
- Add, edit and delete people, each with up to two parents
- Interactive 3D family graph with DAG orientation modes (top-down, bottom-up, radial, free)
- Sortable/searchable people table
- Upcoming and past anniversaries on the home page

## How it works

### Data model

Everything is a graph node. `Person` and `User` embed a common `Node` (UUID, creation/update timestamps, owner). Parenthood is a `PARENT_OF` relationship written on save: a person carries `parent1`/`parent2` UUIDs, and the write query `MERGE`s both the node and its edges in one shot, so the graph stays consistent without a separate relationship-management step.

Neo4j is a natural fit here — a family tree is a DAG, and traversals like "all descendants of X" are the database's native operation rather than a pile of recursive joins.

### Rendering

Server-side rendered `html/template` pages, embedded into the binary with `go:embed` — the compiled binary carries its own templates, no asset directory to deploy.

The graph page feeds nodes and edges to [3d-force-graph](https://github.com/vasturiano/3d-force-graph) (Three.js). Layout is force-directed, with `dagMode` switchable at runtime so the same tree can be read as a hierarchy or as a free-floating network.

### Structure

```
main.go      routes and startup
config/      TOML config loading
db/          Neo4j driver, Cypher queries
sessions/    gorilla/sessions cookie store
types/       Node, Person, User, template Context
views/       HTTP handlers, template rendering
templates/   gohtml pages + layouts (embedded)
```

Handlers behind `views.RequiresLogin` are the whole authorization model.

## Stack

| Component | Tech |
|---|---|
| Language | Go 1.25 |
| HTTP | stdlib `net/http` |
| Templates | stdlib `html/template` + `go:embed` |
| Database | Neo4j (`neo4j-go-driver` v1.8) |
| Sessions | `gorilla/sessions` |
| Password hashing | `golang.org/x/crypto/bcrypt` |
| Config | `pelletier/go-toml` |
| Graph rendering | [3d-force-graph](https://github.com/vasturiano/3d-force-graph) + Three.js |
| UI | Bootstrap 5, DataTables, Select2 (via CDN) |

## Configuration

TOML file, by default at `~/.config/leafnet/config.toml`:

```toml
neo4j_endpoint = "neo4j-db-endpoint"
neo4j_port = "neo4j-db-port"
neo4j_schema = "neo4j-db-schema"
neo4j_username = "neo4j-db-username"
neo4j_password = "neo4j-db-password"
leafnet_port = "leafnet-port"
```

## Running

Needs a reachable Neo4j instance.

```bash
go build
./leafnet -s my-session-key -c config/file/path
```

`-s` can be replaced by the `LSESSION_KEY` environment variable. `-c` defaults to `~/.config/leafnet/config.toml`. The app then listens on `leafnet_port`.

## Known limitations

- Several Cypher queries are built with `fmt.Sprintf` on user-controlled values instead of parameters — injectable. Don't expose this to the internet as is.
- `CSRFToken` exists in the template context but isn't enforced.
- Ownership is recorded on nodes but not checked on read: any logged-in user sees every person in the database.

## License

GPL-3.0 — see [LICENSE](LICENSE).
