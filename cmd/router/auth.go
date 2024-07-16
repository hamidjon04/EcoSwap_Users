package auth

import (
	"ecoswap/api"
	"ecoswap/config"
	"ecoswap/storage/postgres"
	"ecoswap/storage/redis"
	"log"
)

func AuthRun(){
	db, err := postgres.ConnectDB()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	rdb := redis.ConnectRedis()
	defer rdb.Close()

	r := api.Router(db, rdb)
	r.Run(config.Load().USER_ROUTER)
}