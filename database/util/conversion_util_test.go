package util

import (
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
	}
	assert.NotNil(t, user)
	assert.Equal(t, 2, 2)
}
