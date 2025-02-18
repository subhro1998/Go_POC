package app

import (
	"Go_Assignment/auth"
	"Go_Assignment/database/util"
	"Go_Assignment/dto"
	"encoding/json"
	"log"
	"net/http"
)

func Login(writer http.ResponseWriter, requestBody *http.Request) {
	db := util.CreateConnection(false, false)
	if db == nil {
		log.Fatal("No open SQL DB Connections")
		panic(db)
	}
	defer util.CloseDBConnection(db)

	var credentialRequest dto.CredentialRequest
	err := json.NewDecoder(requestBody.Body).Decode(&credentialRequest)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	statusCode, tokenOrErrorMessage := auth.GenerateTokenIfValidUser(db, credentialRequest)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write([]byte(tokenOrErrorMessage))
}

func RefreshToken(writer http.ResponseWriter, requestBody *http.Request) {
	tokenString := requestBody.Header.Get("Authorization")
	auth.RefreshToken(tokenString)
}
