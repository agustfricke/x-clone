package models


type Post struct {
  gorm.Model
	Body        string  
	Image       string  
	UserID      uint
  User        User  
}

