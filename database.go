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

	db_driver := os.Getenv("APP_DB_DRIVER")
	if db_driver == "" {
		db_driver = "mysql"
	}
	db_user := os.Getenv("APP_DB_USER")
	db_pass := os.Getenv("APP_DB_PASS")
	db_location := os.Getenv("APP_DB_LOCATION")
	db_database := os.Getenv("APP_DB_DATABASE")
	db_config := db_user + ":" + db_pass + "@" + db_location + "/" + db_database + "?loc=Local&parseTime=true&charset=utf8"

	var db, db_err = gorm.Open(db_driver, db_config)

	if db_err != nil {
		log.Fatal("Error Database Access")
	}

	//sql log
	db.LogMode(true)

	return db, db_err
}
