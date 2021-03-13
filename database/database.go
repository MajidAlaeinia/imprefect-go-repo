package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/vandario/govms-ipg/app/exceptions"
	"os"
)

func VandarDatabase() (db *gorm.DB) {
	godotenvError := godotenv.Load()
	exceptions.PanicException(godotenvError)
	db, err := gorm.Open("mysql", os.Getenv("DATABASE_VANDAR_USERNAME")+":"+os.Getenv("DATABASE_VANDAR_PASSWORD")+"@/vandar?charset=utf8&parseTime=True&loc=Local")
	exceptions.PanicException(err)
	return db
}
