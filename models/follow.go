package models

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	UserID        uint
	FollowerID    uint 
	FollowingUserID uint 
	Follower      User `gorm:"foreignkey:FollowerID"`
	FollowingUser User `gorm:"foreignkey:FollowingUserID"`
}
