package pg

import "github.com/jackc/pgx/v5/pgxpool"

var db *pgxpool.Pool

func SetPool(pool *pgxpool.Pool) {
	db = pool
}
