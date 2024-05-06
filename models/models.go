package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID uint `gorm:"primarykey" json:"user_id"`
	//gorm.Model
	FistName  string     `gorm:"not null" json:"firstname"`
	LastName  string     `json:"lastname"`
	Email     string     `gorm:"unique;not null" json:"email"`
	Gender    string     `gorm:"check:gender IN ('male','MALE','female','FEMALE','')" json:"gender"`
	Phone     string     `gorm:"not null" json:"phone_no"`
	Password  string     `gorm:"not null" json:"password"`
	Status    string     `gorm:"default:Active" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time  
}

type OTP struct {
	ID uint `gorm:"primarykey" json:"id"`
	//gorm.Model
	Otp    string    `json:"otp"`
	Exp    time.Time //OTP expiry time
	UserID uint      //Foreign key referencing the user model
}

type Admin struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

type Category struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatesAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
}

type Product struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Stock       string    `json:"stock"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatesAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
	CategoryID  uint
	Category    Category
}

// type ProductImage struct {
// 	gorm.Model
// 	ProductID uint
// 	Filename  string `json:"filename"`
// 	URL       string `json:"url"`
// }
