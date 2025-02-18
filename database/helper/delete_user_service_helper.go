package helper

import (
	"Go_Assignment/database/util"
	"Go_Assignment/dto"
	"Go_Assignment/loggerutil"
	"Go_Assignment/model"
	"log"

	"gorm.io/gorm"
)

// Deletes and returns no of users deleted if no user found then returns -1
func DeleteUserIfExists(db *gorm.DB, user model.User) int {
	logMessages := []loggerutil.LogMessage{}
	existingUsers := FetchUsersWithProvidedDetails(db, user)
	if len(existingUsers) == 0 {
		logMsg := loggerutil.LogMessage{Level: "Info",
			Message: "There is no existing user with provided details, will not be able to delete"}
		logMessages = append(logMessages, logMsg)
		// process here and return
		go loggerutil.PostLogMessages(logMessages)
		go loggerutil.LogProcessor()

		log.Println("There is no existing user with provided details, will not be able to delete")
		return 0
	}

	existingUserEntityList := make([]model.User, 0)
	for _, userResponse := range existingUsers {
		userEntity := util.ToDomainDto(dto.UserRequest(userResponse))
		existingUserEntityList = append(existingUserEntityList, userEntity)
	}

	// delete user
	result := db.Unscoped().Delete(&existingUserEntityList)
	//db.Select(clause.Associations).Delete(&user)

	var deleteResult int
	if result.Error != nil {
		logMsg := loggerutil.LogMessage{Level: "Error", Message: "Error occurred while deleting"}
		logMessages = append(logMessages, logMsg)

		log.Println("Error occurred while deleting - ", result.Error)
		deleteResult = -1
	} else {
		logMsg := loggerutil.LogMessage{Level: "Info", Message: "User deleted successfully"}
		logMessages = append(logMessages, logMsg)

		log.Println("User deleted successfully")
		deleteResult = len(existingUsers)
	}

	// first post messages and then write to file
	go loggerutil.PostLogMessages(logMessages)
	go loggerutil.LogProcessor()

	return deleteResult
}
