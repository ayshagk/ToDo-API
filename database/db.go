package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"todoapi/models"
)

var Db *gorm.DB

func InitDB() {
	connStr := "user=postgres host=localhost password=password dbname=toDoApp sslmode=disable"

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
