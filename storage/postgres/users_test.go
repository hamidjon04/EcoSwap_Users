package postgres

import (
	"database/sql"
	pb "ecoswap/genproto/users"
	"log"
	"testing"
)

func Connect()*sql.DB{
	db, err := ConnectDB()
	if err != nil{
		log.Panicln(err)
	}
	return db
}

func TestRegister(t *testing.T){
	db := Connect()
	defer db.Close()

	repo := NewUsersRepo(db)

	req := pb.UserRegister{
		Username: "hamidjon",
		Email: "nuriddinovhamidjon2@gmail.com",
		Password: "$2a$10$4v1Wp.XigUVs1GNxVOGJ7u5AgMlTNvOquggtceDp7Rgpo6JKRThqu",
		FullName: "Nuriddinov Hamidjon",
	}

	err := repo.Register(&req)
	if err != nil {
		t.Error(err)
	}
}