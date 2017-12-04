package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func DbConnect() (*gorm.DB, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var db, db_err = gorm.Open("mysql", os.Getenv("APP_DB_USER")+":"+os.Getenv("APP_DB_PASS")+"@"+os.Getenv("APP_DB_LOCATION")+"/"+os.Getenv("APP_DB_DATABASE")+"?loc=Local")

	if db_err != nil {
		log.Fatal("Error Database Access")
	}

	//sql log
	db.LogMode(true)

	return db, db_err
}
