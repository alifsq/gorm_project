package config

import (
	"database/sql"
	"fmt"
	"time"

)

func InitDB(cfg *Config) *sql.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBCharset,
		cfg.DBParseTime,
		cfg.DBLoc,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed Connect To Mysql: %v", err))
	}

	// Setting Connection Pool (Best Practice Production)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Mastiin beneran nyambung
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed ping MySQL: %v", err))
	}

	fmt.Println("DB connect success")
	return db
}
