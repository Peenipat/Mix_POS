package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
	 "myapp/models"
)
const (
	host = "localhost"
	port = 5432
	user = "myuser"
	password = "mypassword"
	dbname = "mydatabase"
)

var DB *gorm.DB

func ConnectDB(){
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
	db ,err  := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		panic("failed to connect to database")
	} 

	db.AutoMigrate(&models.User{})

	DB = db
}