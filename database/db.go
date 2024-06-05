package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var dbConfig *dbinfo

func Conect() {
	dbConfig = setupDB(os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))

	dsn := dbConfig.getStringPath()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Connection to the database could not be established!\n ", err.Error())
		os.Exit(2)
	}

	log.Println("Connection to the database has been successfully established")

	DB.Logger = logger.Default.LogMode(logger.Info)

	// #TODO: ADD MIGRATION
}
