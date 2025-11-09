package database

import (
	"fmt"
	"log"
	"os"
	"todoapi/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//established a connection to the db using the connection string, opened the connection with gorm and migrated the user and todo tables into postgres
var Db *gorm.DB

func InitDB() {
	//add postgress on render to get database)url  in environmental vars so app can connect to db instead of local host.
connStr := os.Getenv("DATABASE_URL")
if connStr == "" {
    fmt.Println("No DATABASE_URL found! Running locally?")
    connStr = "user=postgres host=localhost password=password dbname=toDoApp sslmode=disable"
}
	//connStr := "user=postgres host=localhost password=password dbname=toDoApp sslmode=disable"

	var err error
	Db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = Db.AutoMigrate(&models.User{}, &models.Todo{})
	if err != nil {
		log.Fatal("error migrating tables", err)
	}

}
