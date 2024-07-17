package auth

import (
	"database/sql"
	"ecoswap/api"
	"ecoswap/config"

	"github.com/redis/go-redis/v9"
)

func AuthRun(db *sql.DB, rdb *redis.Client) {
	r := api.Router(db, rdb)
	r.Run(config.Load().USER_ROUTER)
}
