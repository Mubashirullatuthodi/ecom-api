package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `gorm:"unique;not null" json:"email"`
	Gender    string `gorm:"check:gender IN ('male','MALE','female','FEMALE','')" json:"gender"`
	Phone     string `gorm:"not null" json:"phone_no"`
	Password  string `gorm:"not null" json:"password"`
	Status    string `gorm:"default:Active" json:"status"`
}

type OTP struct {
	ID     uint      `gorm:"primarykey" json:"id"`
	Otp    string    `json:"otp"`
	Email  string    `gorm:"unique" json:"email"`
	Exp    time.Time //OTP expiry time
	UserID uint      //Foreign key referencing the user model
}

type Admin struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64           `json:"price"`
	Quantity    string         `json:"Quantity"`
	ImagePath   pq.StringArray `gorm:"type:text[]" json:"image_path"`
	Status      string         `json:"status"`
	CategoryID  uint           `json:"category_id"`
	Category    Category
}

type Category struct {
	gorm.Model
	Name        string `json:"category_name"`
	Description string `json:"category_description"`
}

type Address struct {
	gorm.Model
	Address  string `json:"address"`
	Town     string `json:"town"`
	District string `json:"district"`
	Pincode  string `json:"pincode"`
	State    string `json:"state"`
	User_ID  uint   `json:"user_id"`
	User     User
}

type Cart struct {
	ID         uint `gorm:"primaryKey"`
	User_ID    uint `json:"user_id"`
	Product_ID uint `json:"product_id"`
	Quantity   uint `json:"quantity"`
	User       User
	Product    Product
}
