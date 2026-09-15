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
	repo := user.NewUserRepository(q, gormDB)

	// Dummy data
	dummyUser := &model.User{
		Name:  "Dummy U",
		Email: "dummy@mail.com",
	}

	ctx := context.Background()
	err = repo.Create(ctx, dummyUser)
	if err != nil {
		log.Fatal(err)
	}

	users, err := repo.GetAllPaginate(ctx, 2, 0)
	if err != nil {
		log.Printf("Gagal mengambil data dari database: %v\n", err)
		return // atau return err tergantung layer
	}

	// 2. Cek jika data memang kosong (opsional, jika butuh log khusus)
	if len(users) == 0 {
		log.Println("Data user kosong")
		return
	}

	// 3. Iterasi slice data user
	for _, u := range users {
		log.Printf("ID: %d, Name: %s, Email: %s\n", u.ID, u.Name, u.Email)
	}

	fmt.Printf("User berhasil dibuat! ID: %d\n", dummyUser.ID)

}
