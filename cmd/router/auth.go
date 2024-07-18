package auth

import (
	"database/sql"
	"ecoswap/api"
	"ecoswap/config"
	"log"

	"github.com/redis/go-redis/v9"
)

func AuthRun(db *sql.DB, rdb *redis.Client) {
	r := api.Router(db, rdb)
	err := r.Run(config.Load().USER_ROUTER)
	if err != nil{
		log.Fatal(err)
	}
}
