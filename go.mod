module github.com/TicketsBot/patreon-db-sync

go 1.25.3

//replace github.com/TicketsBot-cloud/common => ../common

//replace github.com/TicketsBot-cloud/database => ../database

require (
	github.com/TicketsBot-cloud/common v0.0.0-20250509064208-a2d357175463
	github.com/TicketsBot-cloud/database v0.0.0-20250603194547-7b95c33be9d4
	github.com/caarlos0/env/v11 v11.2.2
	github.com/getsentry/sentry-go v0.33.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v4 v4.18.3
	go.uber.org/zap v1.27.0
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
)
