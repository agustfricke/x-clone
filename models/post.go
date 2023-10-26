package models

import "gorm.io/gorm"


type Post struct {
  gorm.Model
	Content string  
	Image   string  
  liked   []User
  repost  []User
	UserID  uint
  User    User  
}

type Comment struct {
  gorm.Model
	Content string  
	Image   string  
  liked   []User
  repost  []User
	UserID  uint
  User    User  
}
