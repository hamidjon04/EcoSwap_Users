package main

import (
	auth "ecoswap/cmd/router"
	"ecoswap/config"
	"ecoswap/genproto/users"
	"ecoswap/service"
	"ecoswap/storage/postgres"
	"ecoswap/storage/redis"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", config.Load().USER_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rdb := redis.ConnectRedis()
	defer rdb.Close()

	userservice := service.NewUsersService(db, rdb)
	service := grpc.NewServer()
	users.RegisterUsersServiceServer(service, userservice)

	go auth.AuthRun(db, rdb)
	time.Sleep(1 * time.Second)
	
	fmt.Printf("Server is listening on port %s\n", config.Load().USER_SERVICE)
	if err = service.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
