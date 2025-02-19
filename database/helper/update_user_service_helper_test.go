package helper

import (
	"Go_Assignment/database/util"
	"Go_Assignment/model"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateUserIfExists_SuccessfulUpdate(t *testing.T) {
	db, mock, err := util.SetUpMockDB()
	assert.NoError(t, err)

	// Fetch existing user mock
	expectedUserQuery := "SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL"
	expectedPrivilegeQuery := "SELECT * FROM `user_privileges` WHERE `user_privileges`.`user_id` = ? AND `user_privileges`.`deleted_at` IS NULL"

	// Set up the expected queries in the mock
	mock.ExpectQuery(expectedUserQuery).WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{"id", "Username", "UserMailId", "UserRole", "Password"}).AddRow(
			1, "John Doe", "email@test.com", "User", "AdminPass1"))
	mock.ExpectQuery(expectedPrivilegeQuery).WillReturnRows(sqlmock.NewRows([]string{"PrivilegeType", "UserID"}).AddRow("All", 1))

	// Update user
	mock.ExpectBegin()
	expectedUpdateUserQuery := "UPDATE `users` SET `id`=?,`updated_at`=?,`user_name`=?,`user_email`=?,`user_role`=?,`user_password`=? WHERE `users`.`deleted_at` IS NULL AND `id` = ?"
	expectedUpdatePrivilegeQuery := "INSERT INTO `user_privileges` (`created_at`,`updated_at`,`deleted_at`,`privilege_type`,`user_id`,`id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`)"

	mock.ExpectExec(expectedUpdateUserQuery).WithArgs(1, sqlmock.AnyArg(), "XYZ", "xyz@email.com", "Admin", "admin_pass", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(expectedUpdatePrivilegeQuery).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sql.NullTime{Time: time.Now(), Valid: false}, "All", 1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// User request
	userPrivilege := model.UserPrivilege{
		UserID:        1,
		Model:         gorm.Model{ID: 2},
		PrivilegeType: "All",
	}

	userReq := model.User{
		Model:      gorm.Model{ID: 1},
		Username:   "XYZ",
		UserMailId: "xyz@email.com",
		Password:   "admin_pass",
		UserRole:   "Admin",
		Privileges: []model.UserPrivilege{userPrivilege},
	}

	// Call update method
	status, statusMessage := UpdateUserIfExists(db, userReq)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There are some unfulfilled expectations: %s", err)
	}

	assert.Equal(t, status, http.StatusAccepted)
	assert.Equal(t, statusMessage, "User details updatedd successsfully")
}

func TestCheckForUpdateInPrivileges_UpdateRequired(t *testing.T) {
	existingUser := model.User{
		Model:      gorm.Model{ID: 1},
		Username:   "UserName",
		UserMailId: "mail@id.com",
		UserRole:   "Admin",
		Password:   "user_password",
		Privileges: []model.UserPrivilege{
			{
				Model:         gorm.Model{ID: 2},
				PrivilegeType: "View",
				UserID:        1,
			},
			{
				Model:         gorm.Model{ID: 3},
				PrivilegeType: "Add",
				UserID:        1,
			},
		},
	}

	updateUserRequest := model.User{
		Model:      gorm.Model{ID: 1},
		Username:   "UserName",
		UserMailId: "mail@id.com",
		UserRole:   "Admin",
		Password:   "user_password",
		Privileges: []model.UserPrivilege{
			{
				Model:         gorm.Model{ID: 2},
				PrivilegeType: "View",
				UserID:        1,
			},
			{
				Model:         gorm.Model{ID: 3},
				PrivilegeType: "Add",
				UserID:        1,
			},
			{
				Model:         gorm.Model{ID: 4},
				PrivilegeType: "Delete",
				UserID:        1,
			},
		},
	}

	privileges, err := checkForUpdateInPrivileges(updateUserRequest, existingUser)
	assert.NotNil(t, privileges)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(privileges))
}

func TestCheckForUpdateInPrivileges_NoUpdate(t *testing.T) {
	existingUser := model.User{
		Model:      gorm.Model{ID: 1},
		Username:   "UserName",
		UserMailId: "mail@id.com",
		UserRole:   "Admin",
		Password:   "user_password",
		Privileges: []model.UserPrivilege{
			{
				Model:         gorm.Model{ID: 2},
				PrivilegeType: "View",
				UserID:        1,
			},
		},
	}

	updateUserRequest := model.User{
		Model:      gorm.Model{ID: 1},
		Username:   "UserName",
		UserMailId: "mail@id.com",
		UserRole:   "Admin",
		Password:   "user_password",
		Privileges: []model.UserPrivilege{
			{
				Model:         gorm.Model{ID: 2},
				PrivilegeType: "View",
				UserID:        1,
			},
		},
	}

	privileges, err := checkForUpdateInPrivileges(updateUserRequest, existingUser)
	assert.Nil(t, err)
	assert.Nil(t, privileges) // No update required scenario
}

func TestCheckForUpdateInPrivileges_NonAdminUser_Failure(t *testing.T) {
	existingUser := model.User{
		Model:      gorm.Model{ID: 1},
		Username:   "UserName",
		UserMailId: "mail@id.com",
		UserRole:   "User",
		Password:   "user_password",
		Privileges: []model.UserPrivilege{
			{
				Model:         gorm.Model{ID: 2},
				PrivilegeType: "View",
				UserID:        1,
			},
		},
	}

	updateUserRequest := model.User{
		Model:      gorm.Model{ID: 1},
		Username:   "UserName",
		UserMailId: "mail@id.com",
		UserRole:   "User",
		Password:   "user_password",
		Privileges: []model.UserPrivilege{
			{
				Model:         gorm.Model{ID: 2},
				PrivilegeType: "View",
				UserID:        1,
			},
			{
				Model:         gorm.Model{ID: 3},
				PrivilegeType: "Add",
				UserID:        1,
			},
			{
				Model:         gorm.Model{ID: 4},
				PrivilegeType: "Delete",
				UserID:        1,
			},
		},
	}

	privileges, err := checkForUpdateInPrivileges(updateUserRequest, existingUser)
	assert.NotNil(t, err)
	assert.Equal(t, "not an admin user, not allowed to update role", err.Error())
	assert.Nil(t, privileges) // No update required scenario
}
