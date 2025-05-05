package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	configuration "github.com/RobsonDevCode/go-profile-service/src/internal/config"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

type Database struct {
	db *sql.DB
}

var (
	instance *Database
	once     sync.Once
)

func NewUserDataBase(config configuration.Config) *sql.DB {
	dbConfig := config.Database
	fmt.Print("getting here")
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Falied to open database connection: %v", err)
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping database after connection has been made: %v", err)
		}
		instance = &Database{
			db: db,
		}
	})

	if instance == nil {
		log.Fatal("Database instance is nil after initialization")
	}
	return instance.db
}

func (d *Database) GetDatebase() *sql.DB {
	return d.db
}

func (d *Database) Close() error {
	return d.db.Close()
}
