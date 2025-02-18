package excelreader

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// User Data structure
type UserData struct {
	Username   string
	UserMailId string
	UserRole   string
	Privileges []string
	Password   string
}

func ReadExcel(filepath string) []UserData {
	// Opening the Excel file
	file, errReadFile := excelize.OpenFile(filepath)
	checkError(errReadFile)

	// Read data from a all cells
	allRows, errGetRows := file.GetRows("Facility_Data") // Read from Facility_Data sheet
	defer file.Close()                                   // Closes the opened file once main method exists
	checkError(errGetRows)

	// Extract data from all rows and map them to UserData slice
	excelData := []UserData{}
	for index, row := range allRows {
		// Skip the header row
		if index == 0 {
			continue
		}

		var user UserData
		if len(row) <= 5 {
			user.Username = row[0]
			user.UserMailId = row[1]
			user.UserRole = row[2]
			user.Password = row[4]

			privileges := extractPrivileges(row, index)
			user.Privileges = privileges // Set Privileges in UserData
		}

		// Append in the User List
		excelData = append(excelData, user)
	}
	return excelData
}

// Extract and store privileges in a slice
func extractPrivileges(row []string, index int) []string {
	if row[3] == "" {
		fmt.Printf("There is no privileges attached to %v user\n", index+1)
		return nil
	}

	privilegeSplitValues := strings.Split(row[3], ",") // Split based on ','
	// Trim any extra spaces
	for i := range privilegeSplitValues {
		privilegeSplitValues[i] = strings.TrimSpace(privilegeSplitValues[i])
	}
	return privilegeSplitValues
}

// Check for any error
func checkError(err error) {
	if err != nil {
		fmt.Println("Error occurred during execution")
		panic(err)
	}
}

func PrintAllUserDetails(users []UserData) {
	for _, user := range users {
		fmt.Printf("User Name: %v, User Email ID: %v, UserRole: %v, User Privileges: %v\n", user.Username, user.UserMailId,
			user.UserRole, user.Privileges)
	}
}
