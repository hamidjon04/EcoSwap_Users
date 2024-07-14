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

	router := api.Router(db)
	router.Run(config.Load().USER_SERVICE)
}