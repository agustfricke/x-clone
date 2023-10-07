package models

import "gorm.io/gorm"


type Post struct {
  gorm.Model
	Body        string  
	Image       string  
	UserID      uint
  User        User  
}

