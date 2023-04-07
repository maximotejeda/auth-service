package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	DB = NewDB()
	DB.AutoMigrate(&User{}, &UserInfo{}, &RecoverAccount{})
}

// Connection to sqlite
func NewDB() *gorm.DB {
	conn := createConnSTR()
	db, err := gorm.Open(sqlite.Open(conn), &gorm.Config{})
	if err != nil {
		panic("Failed to open DB")
	}
	return db
}

func createConnSTR() string {
	dbPath := os.Getenv("DBDIR")
	dbName := os.Getenv("DBFILE")
	mode := os.Getenv("MODE")
	cache := os.Getenv("CACHE")
	err := os.Mkdir("db", 0777)
	if err != nil {
		fmt.Printf("%v", err)
	}
	return fmt.Sprintf("%s/%s?mode=%s&cache=%s", dbPath, dbName, mode, cache)
}
