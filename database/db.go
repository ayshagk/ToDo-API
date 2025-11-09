package database

import (
	"log"
	"os"
	"todoapi/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//established a connection to the db using the connection string, opened the connection with gorm and migrated the user and todo tables into postgres
var Db *gorm.DB

func InitDB() {
	//had to add this for it to work on render
	connStr := os.Getenv("DATABASE_URL")
if connStr == "" {
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
