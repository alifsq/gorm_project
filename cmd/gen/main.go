package main

import (
	"go-gorm/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB := config.InitDB(cfg)
	defer sqlDB.Close()

	gormDB, err := config.InitGORM(sqlDB)
	if err != nil {
		log.Fatal("Failed Iinit GORM", err)
	}
	config.RunGenerator(
		gormDB,
		"./internal/database/query",
		"./internal/database/model",
	)
}
