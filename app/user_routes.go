package app

import (
	"Go_Assignment/database/helper"
	"Go_Assignment/database/util"
	"Go_Assignment/dto"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func FetchAllUsers(writer http.ResponseWriter, requestBody *http.Request) {
	db := util.CreateConnection(false, true)
	if db == nil {
		log.Fatal("No open SQL DB Connections")
		panic(db)
	}
	defer util.CloseDBConnection(db)

	allUsers := helper.FetchAllUsers(db)
	var responseCode int
	if allUsers != nil || len(allUsers) > 0 {
		responseCode = http.StatusFound
	} else {
		responseCode = http.StatusNotFound
	}

	usersJson, jsonErr := json.Marshal(allUsers)
	if jsonErr != nil {
		log.Fatal("Not able to convert User response to JSON obj")
		panic(jsonErr)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(responseCode)

	// Encode response to JSON payload as ResponseBody
	writer.Write(usersJson)
}

func FetchSpecificUser(writer http.ResponseWriter, requestBody *http.Request) {
	db := util.CreateConnection(false, true)
	if db == nil {
		log.Fatal("No open SQL DB Connections")
		panic(db)
	}
	defer util.CloseDBConnection(db)

	user := dto.UserRequest{}
	err := json.NewDecoder(requestBody.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	userEntity := util.ToDomainDto(user)
	userResponseList := helper.FetchUsersWithProvidedDetails(db, userEntity)

	var responseCode int
	if userResponseList != nil || len(userResponseList) > 0 {
		responseCode = http.StatusFound
	} else {
		responseCode = http.StatusNotFound
	}

	usersJson, jsonErr := json.Marshal(userResponseList)
	if jsonErr != nil {
		log.Fatal("Not able to convert User response to JSON obj")
		panic(jsonErr)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(responseCode)

	// Encode response to JSON payload as ResponseBody
	writer.Write(usersJson)
}

func SaveUser(writer http.ResponseWriter, requestBody *http.Request) {
	db := util.CreateConnection(false, true)
	if db == nil {
		log.Fatal("No open SQL DB Connections")
		panic(db)
	}
	defer util.CloseDBConnection(db)

	users := []dto.UserRequest{}
	err := json.NewDecoder(requestBody.Body).Decode(&users)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if len(users) <= 0 {
		http.Error(writer, "No Request body received", http.StatusBadRequest)
		return
	}

	// Convert user to interface for using generic type
	userDataAbstract := make([]interface{}, len(users))
	for i, user := range users {
		userDataAbstract[i] = user
	}

	status, saveResponseMessage := helper.SaveUserDetails(db, userDataAbstract)
	saveResponseJson, jsonErr := json.Marshal(saveResponseMessage)
	if jsonErr != nil {
		log.Fatal("Not able to convert User response to JSON obj")
		panic(jsonErr)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	// Encode response to JSON payload as ResponseBody
	writer.Write(saveResponseJson)
}

// Checks if user exists and if it exists then updates the user
func UpdateUser(writer http.ResponseWriter, requestBody *http.Request) {
	db := util.CreateConnection(false, true)
	if db == nil {
		log.Fatal("No open SQL DB Connections")
		panic(db)
	}
	defer util.CloseDBConnection(db)

	userRequest := dto.UserRequest{}
	err := json.NewDecoder(requestBody.Body).Decode(&userRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	userEntity := util.ToDomainDto(userRequest)
	status, updateResponseMessage := helper.UpdateUserIfExists(db, userEntity)
	updateResponseJson, jsonErr := json.Marshal(updateResponseMessage)
	if jsonErr != nil {
		log.Fatal("Not able to convert User response to JSON obj")
		panic(jsonErr)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	// Encode response to JSON payload as ResponseBody
	writer.Write(updateResponseJson)
}

func DeleteUser(writer http.ResponseWriter, requestBody *http.Request) {
	db := util.CreateConnection(false, true)
	if db == nil {
		log.Fatal("No open SQL DB Connections")
		panic(db)
	}
	defer util.CloseDBConnection(db)

	users := []dto.UserRequest{}
	err := json.NewDecoder(requestBody.Body).Decode(&users)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if len(users) != 1 {
		http.Error(writer, "No of User in request body should be only 1", http.StatusBadRequest)
		return
	}

	userEntity := util.ToDomainDto(users[0])                       // Convert request into Entity
	numOfUsersDeleted := helper.DeleteUserIfExists(db, userEntity) // Delete user/ users if exists based on provided conditions

	var responseCode int
	var responseStr string
	if numOfUsersDeleted > 0 {
		responseCode = http.StatusAccepted
		responseStr = strconv.Itoa(numOfUsersDeleted) + " users successfully Deleted."
	} else {
		responseCode = http.StatusNotFound
		responseStr = "No user found with provided details to be deleted"
	}

	deleteResponseJson, jsonErr := json.Marshal(responseStr)
	if jsonErr != nil {
		log.Fatal("Not able to convert User response to JSON obj")
		panic(jsonErr)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(responseCode)

	// Encode response to JSON payload as ResponseBody
	writer.Write(deleteResponseJson)
}
