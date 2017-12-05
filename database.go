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

	// Read Env Dababase params
	db_user := os.Getenv("APP_DB_USER")
	db_pass := os.Getenv("APP_DB_PASS")
	db_location := os.Getenv("APP_DB_LOCATION")
	db_database := os.Getenv("APP_DB_DATABASE")

	// Database Config
	db_config := db_user + ":" + db_pass + "@" + db_location + "/" + db_database + "?charset=utf8&parseTime=True&loc=Local"

	var db, db_err = gorm.Open("mysql", db_config)

	if db_err != nil {
		log.Fatal("Error Database Access")
	}

	//sql log
	db.LogMode(true)

	return db, db_err
}
