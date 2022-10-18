package pgdatabase

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func InitDB(ctx context.Context, url string) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, "postgres", url)
}
