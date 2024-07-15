package main

import (
	"ecoswap/api"
	"ecoswap/config"
	"ecoswap/storage/postgres"
	"log"
)

func main(){
	db, err := postgres.ConnectDB()
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	r := api.Router(db)
	r.Run(config.Load().USER_ROUTER)
}