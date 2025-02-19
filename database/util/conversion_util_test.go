package util

import (
	"Go_Assignment/dto"
	"Go_Assignment/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestToResponseDto(t *testing.T) {
	user := model.User{
		Username:   "Test_Name",
		UserMailId: "Test_email@deloitte.com",
		UserRole:   "Test_Role",
		Model:      gorm.Model{ID: 1},
		Password:   "Test_Pass",
		Privileges: []model.UserPrivilege{
			{
				Model:         gorm.Model{ID: 2},
				PrivilegeType: "Test_Privilege_Type1",
				UserID:        1,
			},
			{
				Model:         gorm.Model{ID: 3},
				PrivilegeType: "Test_Privilege_Type2",
				UserID:        1,
			},
		},
	}

	// Call actual method
	response := ToResponseDto(user)

	// compare the returned result
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Privileges))
	assert.Equal(t, "Test_Role", response.Role)
	assert.Equal(t, "Test_email@deloitte.com", response.Email)
}

func TestToDomainDto(t *testing.T) {
	userRequest := dto.UserRequest{
		Name:     "Test_Name",
		Email:    "Test_email@deloitte.com",
		Role:     "Test_Role",
		UserId:   1,
		Password: "Test_Pass",
		Privileges: []dto.UserPrivilegeResponse{
			{
				PrivilegeId:   2,
				PrivilegeType: "Test_Privilege_Type1",
				UserId:        1,
			},
			{
				PrivilegeId:   3,
				PrivilegeType: "Test_Privilege_Type2",
				UserId:        1,
			},
			{
				PrivilegeId:   4,
				PrivilegeType: "Test_Privilege_Type3",
				UserId:        1,
			},
		},
	}

	// Call the method
	response := ToDomainDto(userRequest)

	// Compare the response
	assert.NotNil(t, response)
	assert.Equal(t, 3, len(response.Privileges))
	assert.Equal(t, "Test_Pass", response.Password)
	assert.Equal(t, "Test_email@deloitte.com", response.UserMailId)
}
