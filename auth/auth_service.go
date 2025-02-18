package auth

import (
	"Go_Assignment/database/helper"
	"Go_Assignment/dto"
	"Go_Assignment/model"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

var jwtKey = []byte("secret_key")

// Returns HTTP Status code, Error message if any else tokenString and expirationTime
func GenerateTokenIfValidUser(db *gorm.DB, loginRequest dto.CredentialRequest) (int, string) {
	isValidReq, errMessage := validateLoginRequest(loginRequest)
	if !isValidReq {
		log.Println("In valid login request, necessary parameters are missing")
		return http.StatusBadRequest, errMessage
	}

	// Check the count
	userToSearch := model.User{
		Model:      gorm.Model{ID: uint(loginRequest.UserID)},
		UserMailId: loginRequest.UserEmailId,
		Password:   loginRequest.Password,
	}
	users := helper.FetchUsersWithProvidedDetails(db, userToSearch)

	if len(users) == 1 {
		log.Println("Provided User details is correct.")
	} else {
		log.Println("Error :: In valid for login request or credential")
		return http.StatusUnauthorized, "Error :: In valid for login credential provided"
	}

	// Now create JWT Auth code
	expirationTime := time.Now().Add(time.Minute * 20)
	claims := &dto.Claims{
		UserID:      loginRequest.UserID,
		UserEmailId: loginRequest.UserEmailId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil || len(tokenString) <= 0 {
		log.Println("Not able to generate JWT token")
		return http.StatusInternalServerError, "Not able to generate JWT token"
	}

	return http.StatusOK, tokenString
}

func validateLoginRequest(loginRequest dto.CredentialRequest) (bool, string) {
	message := "Valid request body"
	isValidRequest := true

	if loginRequest.UserID <= 0 && len(loginRequest.UserEmailId) <= 0 {
		message = "Either UserId or Email Id is required for logging in, not present in the request."
		log.Println(message)
		isValidRequest = false
	}

	if len(loginRequest.Password) <= 0 {
		message = "Password cannot be blank for login."
		log.Println(message)
		isValidRequest = false
	}

	return isValidRequest, message
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = getTokenFromHeader(tokenString) // Remove the bearer at the start and get the actual token
		if len(tokenString) <= 0 {
			http.Error(w, "No token provided in Header", http.StatusUnauthorized)
			return
		}

		claims := &dto.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RefreshToken(tokenString string) (int, string) {
	tokenString = getTokenFromHeader(tokenString) // Remove the bearer at the start and get the actual token
	if len(tokenString) <= 0 {
		log.Println("No token provided in Header for refreshing")
		return http.StatusUnauthorized, "No token provided in Header"
	}

	claims := &dto.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || len(tokenString) <= 0 || !token.Valid {
		log.Println("Not able to Parse JWT token")
		return http.StatusUnauthorized, "Not able to parse JWT token"
	}

	if time.Until(time.Unix(claims.ExpiresAt, 0)) > time.Second*30 {
		log.Println("Previous JWT token not yet expired")
		return http.StatusBadRequest, "Previous JWT token not yet expired"
	}

	// Now create JWT Auth code
	expirationTime := time.Now().Add(time.Minute * 30)
	claims.ExpiresAt = expirationTime.Unix()

	// Generate JWT token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString(jwtKey)

	if err != nil || len(newTokenString) <= 0 {
		log.Println("Not able to Generate Refresh JWT token")
		return http.StatusInternalServerError, "Not able to Generate Refresh JWT token"
	}

	return http.StatusOK, newTokenString
}

func getTokenFromHeader(tokenFromHeader string) string {
	splitToken := strings.Split(tokenFromHeader, "Bearer")
	if len(splitToken) == 2 {
		return strings.TrimSpace(splitToken[1])
	}
	return ""
}
