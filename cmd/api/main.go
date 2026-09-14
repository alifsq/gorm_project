package main

import (
	"context"
	"fmt"
	"go-gorm/config"
	"go-gorm/internal/database/model"
	"go-gorm/internal/database/query"
	"go-gorm/internal/modules/user"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	} else {
		log.Print("Init Load env success")
	}

	sqlDB := config.InitDB(cfg)
	defer sqlDB.Close()

	gormDB, err := config.InitGORM(sqlDB)
	if err != nil {
		log.Fatal("Init Gorm Failed : ", err)
	}
	// DI Repository
	q := query.Use(gormDB)
	repo := user.NewUserRepository(q)

	// Dummy data
	dummyUser := &model.User{
		Name:  "Dummy User",
		Email: "dummy@example.com",
	}

	err = repo.Create(context.Background(), dummyUser)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User berhasil dibuat! ID: %d\n", dummyUser.ID)

}
