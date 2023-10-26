package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string
	Email        string
	Password     string
	Verified     bool
	SocialID     string
	Avatar       string
	OtpEnabled   bool
	OtpVerified  bool
	OtpSecret    string
	OtpAuthURL   string

	Posts     []Post
	Followers []Follow
	Following []Follow
}



type GoogleResponse struct {
	ID       string 
	Email    string 
	Verified bool   
	Picture  string 
}

type GitHubResponse struct {
	ID       int
	Email    string 
  Name    string 
	AvatarURL string 
}

