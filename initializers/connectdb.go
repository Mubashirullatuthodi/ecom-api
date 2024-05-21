package initializers

import (
	"log"
	"os"

	"github.com/mubashir/e-commerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(){
	DSN:=os.Getenv("DSN")

	var err error

	dsn := DSN

	DB,err=gorm.Open(postgres.Open(dsn),&gorm.Config{})

	if err != nil{
		log.Fatal("error connecting to database")
	}

	DB.AutoMigrate(&models.User{},&models.Admin{},&models.OTP{},&models.Category{},&models.Product{},&models.Address{},&models.Cart{})
}