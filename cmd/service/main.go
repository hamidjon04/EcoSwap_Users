package main

import (
	"ecoswap/api"
	auth "ecoswap/cmd/router"
	"ecoswap/config"
	"ecoswap/storage/postgres"
	"ecoswap/storage/redis"
	"log"
	"time"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rdb := redis.ConnectRedis()
	defer rdb.Close()

	router := api.Router(db, rdb)
	go auth.AuthRun(db, rdb)
	time.Sleep(1 * time.Second)
	router.Run(config.Load().USER_SERVICE)
}


