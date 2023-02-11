package db

import (
	"fmt"
	"log"

	"github.com/roskeys/app/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() bool {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", utils.DB_USERNAME, utils.DB_PASSWORD, utils.DB_DATABASE)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connect to database")
	}
	DB = db
	return true
}

var _ = InitDatabase()
