package models

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	UserID   uint 
	Follower User 
	FollowingUser User 
}
