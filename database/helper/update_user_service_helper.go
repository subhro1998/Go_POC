package helper

import (
	"Go_Assignment/database/util"
	"Go_Assignment/dto"
	"Go_Assignment/loggerutil"
	"Go_Assignment/model"
	"errors"
	"log"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func UpdateUserIfExists(db *gorm.DB, userDetailsToUpdate model.User) (int, string) {
	logMessages := []loggerutil.LogMessage{}
	logMsg := loggerutil.LogMessage{Level: "Info", Message: "Inside Service call to Update user with provided details"}
	logMessages = append(logMessages, logMsg)
	message, statusCode := "", 0

	existingUsers := FetchUsersWithProvidedDetails(db, model.User{Model: gorm.Model{ID: userDetailsToUpdate.ID}})
	if len(existingUsers) != 1 {
		if len(existingUsers) == 0 {
			message = "No User found in DB with provided details"
			logMsg = loggerutil.LogMessage{Level: "Info", Message: message}
			statusCode = http.StatusNotFound
		} else if len(existingUsers) > 1 {
			message = "More than 1 user is found with provided details"
			logMsg = loggerutil.LogMessage{Level: "Info", Message: message}
			statusCode = http.StatusNotAcceptable
		}
		logMessages = append(logMessages, logMsg)
		go loggerutil.PostLogMessages(logMessages)
		go loggerutil.LogProcessor()

		return statusCode, message
	}

	// convert the Request to Domain entity and check for field Update and then Update those fields
	existingUser := util.ToDomainDto(dto.UserRequest(existingUsers[0]))
	updatedPrivileges, updateErr := checkForUpdateInPrivileges(userDetailsToUpdate, existingUser)
	if updateErr != nil {
		return http.StatusNotAcceptable, updateErr.Error()
	}

	if updatedPrivileges != nil {
		logMsg = loggerutil.LogMessage{Level: "Info", Message: "Going to replace current Privileges to new Privileges"}
		logMessages = append(logMessages, logMsg)
		log.Println("Replacing current Privileges to new Privileges with : ", updatedPrivileges)

		// Replaceing privileges
		privilegeError := db.Model(&userDetailsToUpdate).Association("Privileges").Replace(updatedPrivileges)
		if privilegeError != nil {
			logMsg = loggerutil.LogMessage{Level: "Error", Message: "Error occurred while replacing new Privileges, will be a partial save case"}
			logMessages = append(logMessages, logMsg)
		}
	}

	logMsg = loggerutil.LogMessage{Level: "Info", Message: "Going to update users table deatils"}
	logMessages = append(logMessages, logMsg)
	result := db.Model(&userDetailsToUpdate).Updates(userDetailsToUpdate) // Update existing user

	if result.Error != nil {
		statusCode = http.StatusBadRequest
		message = "Some error occurred while saving the Updated User details"
		logMsg = loggerutil.LogMessage{Level: "Error", Message: message}
	} else {
		statusCode = http.StatusAccepted
		message = "User details updatedd successsfully"
		logMsg = loggerutil.LogMessage{Level: "Info", Message: message}
	}

	logMessages = append(logMessages, logMsg)
	go loggerutil.PostLogMessages(logMessages)
	go loggerutil.LogProcessor()

	return statusCode, message
}

func checkForUpdateInPrivileges(userDetailsToUpdate model.User, existingUserDetails model.User) ([]model.UserPrivilege, error) {
	// check if there is any update in Privileges
	logMessages := []loggerutil.LogMessage{}
	logMsg := loggerutil.LogMessage{Level: "Info", Message: "Checking for update in provided user details"}
	logMessages = append(logMessages, logMsg)

	updateRequired := false
	if len(userDetailsToUpdate.Privileges) != len(existingUserDetails.Privileges) {
		logMsg := loggerutil.LogMessage{Level: "Info", Message: "No of Privileges is different, need to update"}
		logMessages = append(logMessages, logMsg)
		updateRequired = true
	} else {
		privilegeMap := make(map[string]bool)
		for _, privilegeFromUpdateReq := range userDetailsToUpdate.Privileges {
			privilegeMap[privilegeFromUpdateReq.PrivilegeType] = true
		}

		for _, existingPrivilege := range existingUserDetails.Privileges {
			if !privilegeMap[existingPrivilege.PrivilegeType] {
				logMsg := loggerutil.LogMessage{Level: "Info", Message: "There is a update in existing Privileges, need to update"}
				logMessages = append(logMessages, logMsg)
				updateRequired = true
			}
		}
	}

	go loggerutil.PostLogMessages(logMessages)
	go loggerutil.LogProcessor()

	if updateRequired {
		if strings.EqualFold(existingUserDetails.UserRole, "admin") {
			return userDetailsToUpdate.Privileges, nil
		} else {
			return nil, errors.New("not an admin user, not allowed to update role")
		}
	} else {
		return nil, nil
	}
}
