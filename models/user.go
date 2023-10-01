package models

import "gorm.io/gorm"

type User struct {
  gorm.Model
	Name          string  
	Email         string  
	Password      string 
	Verified      bool    
	SocialID      string 
	Avatar        string  
	Otp_enabled   bool   
	Otp_verified  bool  
	Otp_secret    string
	Otp_auth_url  string

	Posts       []Post
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

type SignUpInput struct {
	Name            string 
	Email           string 
	Password        string 
}

type SignInInput struct {
	Email     string 
	Password  string 
  Token     string
}

type OTPInput struct {
	Token   string 
}
