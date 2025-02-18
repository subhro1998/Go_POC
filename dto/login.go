package dto

import "github.com/dgrijalva/jwt-go"

type CredentialRequest struct {
	UserID      int    `json:"user_id"`
	UserEmailId string `json:"user_email"`
	Password    string `json:"password"`
}

type Claims struct {
	UserID      int    `json:"user_id"`
	UserEmailId string `json:"user_email"`
	jwt.StandardClaims
}
