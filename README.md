# ecom-go

A Go e-commerce API with MySQL, Gorilla Mux routing, and SQL migrations. Features are split by domain under `service/`, with shared types and JSON helpers at the top level.

## Folder architecture

```
cmd
 ├── main.go                 # app entrypoint (DB + API server)
 ├── api
 │    └── api.go             # HTTP server setup, route mounting
 └── migrate
      ├── main.go            # migration runner
      └── migrations/        # up/down SQL migrations

config
 └── env.go                  # environment / config loading

db
 └── db.go                   # MySQL connection setup

service
 ├── auth/                   # password hashing / auth helpers
 └── user/
      ├── routes.go          # routes + handlers/controllers
      └── store.go           # DB interactions

types/                       # shared request/response & domain types
utils/                       # WriteJSON / ParseJSON / Validate
```

| Path | Role |
|------|------|
| `cmd/` | Binaries — API process and migrate CLI |
| `cmd/api/` | Wires the router and mounts service handlers |
| `cmd/migrate/` | Applies SQL migrations |
| `config/` | Loads env vars (DB credentials, etc.) |
| `db/` | Database driver / connection helpers |
| `service/<domain>/` | Feature modules: HTTP layer in `routes.go`, persistence in `store.go` |
| `types/` | Shared structs used across services |
| `utils/` | JSON write/parse and validation helpers |
