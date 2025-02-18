package helper

import (
	"Go_Assignment/database/util"
	"Go_Assignment/dto"
	"Go_Assignment/excelreader"
	"Go_Assignment/loggerutil"
	"Go_Assignment/model"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func convertExcelUserDetailsToEntity(userDataList []excelreader.UserData) []model.User {
	userList := []model.User{}
	for _, userData := range userDataList {
		user := model.User{}
		user.Username = userData.Username
		user.UserRole = userData.UserRole
		user.UserMailId = userData.UserMailId
		user.Password = userData.Password

		// create Privilege
		userPrivilegeList := []model.UserPrivilege{}
		for _, privilege := range userData.Privileges {
			userPrivilege := model.UserPrivilege{PrivilegeType: privilege}
			userPrivilegeList = append(userPrivilegeList, userPrivilege)
		}
		user.Privileges = userPrivilegeList // Set user privileges
		userList = append(userList, user)
	}

	return userList
}

// func SaveUserDetails(userDataList []excelreader.UserData) {
func SaveUserDetails(db *gorm.DB, userDataAbstract []interface{}) (int, string) {
	logMessages := []loggerutil.LogMessage{}
	logMsg := loggerutil.LogMessage{Level: "Info", Message: "Inside Service call to save user with provided details"}
	logMessages = append(logMessages, logMsg)

	if len(userDataAbstract) <= 0 {
		logMsg = loggerutil.LogMessage{Level: "Warn", Message: "No user data provided in request body"}
		logMessages = append(logMessages, logMsg)

		go loggerutil.PostLogMessages(logMessages)
		go loggerutil.LogProcessor()

		return http.StatusBadRequest, "No data in Request"
	}

	userEntityList := []model.User{}
	checkDBEntity := CheckStructDataType(userDataAbstract)
	switch checkDBEntity {

	case EXCEL_MODEL: // Save data into DB from Excel in bulk
		logMsg = loggerutil.LogMessage{Level: "Info", Message: "Request body is of Excel Model"}
		logMessages = append(logMessages, logMsg)
		excelUserDataList := []excelreader.UserData{}
		for _, excelUserData := range userDataAbstract {
			userData := excelUserData.(excelreader.UserData)
			excelUserDataList = append(excelUserDataList, userData)
		}
		userEntityList = convertExcelUserDetailsToEntity(excelUserDataList)

	case REQUEST_MODEL: // Save data into DB from JSON request
		logMsg = loggerutil.LogMessage{Level: "Info", Message: "Request body is of DB Model"}
		logMessages = append(logMessages, logMsg)
		for _, userReq := range userDataAbstract {
			userData := userReq.(dto.UserRequest)
			userEntityPostConversion := util.ToDomainDto(userData)
			userEntityList = append(userEntityList, userEntityPostConversion)
		}
	}

	logMsg = loggerutil.LogMessage{Level: "Info", Message: "Request converted to domain dto, going to save data"}
	logMessages = append(logMessages, logMsg)

	// Save User in DB
	//userEntityList := convertUserDetailsToEntity(userDataList)
	result := db.Create(userEntityList)
	var statusMessage string
	var status int
	// Check for errors
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "1062") || strings.Contains(result.Error.Error(), "Duplicate entry") {
			statusMessage = "Email already exists."
			logMsg = loggerutil.LogMessage{Level: "Error", Message: statusMessage}
			logMessages = append(logMessages, logMsg)
		} else {
			log.Println("Failed to save user because of error : ", result.Error)
			statusMessage = "Not able to save provided " + strconv.Itoa(len(userEntityList)) + " users due to an error."
			logMsg = loggerutil.LogMessage{Level: "Error", Message: statusMessage}
			logMessages = append(logMessages, logMsg)
		}
		status = http.StatusBadRequest
	} else {
		status = http.StatusAccepted
		statusMessage = "Provided " + strconv.Itoa(len(userEntityList)) + " User is saved in DB"
		logMsg = loggerutil.LogMessage{Level: "Info", Message: statusMessage}
		logMessages = append(logMessages, logMsg)
	}

	go loggerutil.PostLogMessages(logMessages)
	go loggerutil.LogProcessor()

	return status, statusMessage
}

func CheckStructDataType(t []interface{}) string {
	switch t[0].(type) {
	case excelreader.UserData:
		return EXCEL_MODEL
	case dto.UserRequest:
		return REQUEST_MODEL
	case model.User:
		return DB_ENTITY_MODEL
	default:
		return "NON_EXISTING_MODEL"
	}
}
