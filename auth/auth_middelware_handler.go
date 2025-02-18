package auth

// import (
// 	"Go_Assignment/dto"
// 	"net/http"

// 	"github.com/dgrijalva/jwt-go"
// )

// var jwtKey = []byte("secret_key")

// func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := r.Header.Get("Authorization")
// 		if len(tokenString) <= 0 {
// 		}

// 		claims := &dto.Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }
