package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	username = os.Getenv("ADMINUSERDB")
	password = os.Getenv("ADMINPASSDB")
	email    = os.Getenv("ADMINEMAILDB")
	addAdmin = os.Getenv("ADDADMINUSER")
)

func init() {
	DB = NewDB()
	DB.AutoMigrate(&User{}, &UserInfo{}, &RecoverAccount{})
	// Create admin user on first run of database
	if addAdmin != "" {
		addAdminUser()
	}

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

func addAdminUser() {
	if username != "" && password != "" && email != "" {
		passwd, _ := hashPassword(password)
		DB.Create(&User{
			Active: true,
			Rol:    "admin",
			UserLogin: UserLogin{
				Username: username,
				Password: passwd,
				Email:    email,
			},
			RecoverAccount: RecoverAccount{
				RecoverCount: 0,
				Type:         "email",
				Recover:      0,
			},
			UserInfo: UserInfo{
				FirstName: "admin",
				LastName:  "admin",
			},
		})
	}
}
