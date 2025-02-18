package helper

import (
	"Go_Assignment/database/util"
	"Go_Assignment/dto"
	"Go_Assignment/loggerutil"
	"Go_Assignment/model"
	"log"
	"strconv"

	"gorm.io/gorm"
)

func FetchAllUsers(db *gorm.DB) []dto.UserResponse {
	logMessages := []loggerutil.LogMessage{}
	logMsg := loggerutil.LogMessage{Level: "Info", Message: "Inside Service call to fetch all users"}
	logMessages = append(logMessages, logMsg)

	if db == nil {
		log.Fatal("No open SQL DB Connections")
		panic(db)
	}

	var users []model.User
	result := db.Preload("Privileges").Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
		panic(result.Error)
	}

	logMsg = loggerutil.LogMessage{Level: "Info", Message: "No of Users fetched is " + strconv.Itoa(len(users))}
	logMessages = append(logMessages, logMsg)

	log.Println("No of Users fetched is ", len(users))
	userResponseList := []dto.UserResponse{}
	for _, userDomain := range users {
		userResponse := util.ToResponseDto(userDomain)
		userResponseList = append(userResponseList, userResponse)
	}

	logMsg = loggerutil.LogMessage{Level: "Info", Message: "User is converted to UserResponse"}
	logMessages = append(logMessages, logMsg)

	go loggerutil.PostLogMessages(logMessages)
	go loggerutil.LogProcessor()

	return userResponseList
}

func FetchUsersWithProvidedDetails(db *gorm.DB, userToSearch model.User) []dto.UserResponse {
	logMessages := []loggerutil.LogMessage{}
	logMsg := loggerutil.LogMessage{Level: "Info", Message: "Inside Service call to fetch speciic user with provided details"}
	logMessages = append(logMessages, logMsg)

	var users []model.User
	result := db.Preload("Privileges").Where(&userToSearch).Find(&users)
	if result.Error != nil {
		logMsg = loggerutil.LogMessage{Level: "Error", Message: "Some Error occurred during DB call"}
		logMessages = append(logMessages, logMsg)

		log.Println("Some Error occurred during DB call and the error is : ", result.Error)
		//panic(result.Error)
	}

	log.Println("No of Users fetched is ", len(users))
	userResponseList := []dto.UserResponse{}
	for _, userDomain := range users {
		userResponse := util.ToResponseDto(userDomain)
		userResponseList = append(userResponseList, userResponse)
	}

	logMsg = loggerutil.LogMessage{Level: "Info",
		Message: "Fetched no of User with prived details is " + strconv.Itoa(len(userResponseList))}
	logMessages = append(logMessages, logMsg)

	go loggerutil.PostLogMessages(logMessages)
	go loggerutil.LogProcessor()

	return userResponseList
}
