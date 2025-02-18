package helper

import (
	"Go_Assignment/database/util"
	"Go_Assignment/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteUserIfExists_SuccessfullyDeleted(t *testing.T) {
	db, mock, err := util.SetUpMockDB()
	assert.NoError(t, err)

	// Set up the expected queries in the mock
	expectedUserQuery := "SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL"
	expectedPrivilegeQuery := "SELECT * FROM `user_privileges` WHERE `user_privileges`.`user_id` = ? AND `user_privileges`.`deleted_at` IS NULL"

	// Set up the expected queries in the mock
	mock.ExpectQuery(expectedUserQuery).WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{"id", "Username", "UserMailId", "UserRole", "Password"}).AddRow(
			1, "John Doe", "email@test.com", "User", "AdminPass1"))
	mock.ExpectQuery(expectedPrivilegeQuery).WillReturnRows(sqlmock.NewRows([]string{"PrivilegeType", "UserID"}).AddRow("All", 1))

	// Delete query
	mock.ExpectBegin()
	// Set up the expected queries in the mock
	expectedDeleteUserQuery := "DELETE FROM `users` WHERE `users`.`id` = ?"
	mock.ExpectExec(expectedDeleteUserQuery).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the function
	usersDeleted := DeleteUserIfExists(db, model.User{Model: gorm.Model{ID: 1}})

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, 1, usersDeleted)
}

func TestDeleteUserIfExists_AdminUserNotRemovable(t *testing.T) {
	db, mock, err := util.SetUpMockDB()
	assert.NoError(t, err)

	// Set up the expected queries in the mock
	expectedUserQuery := "SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL"
	expectedPrivilegeQuery := "SELECT * FROM `user_privileges` WHERE `user_privileges`.`user_id` = ? AND `user_privileges`.`deleted_at` IS NULL"

	// Set up the expected queries in the mock
	mock.ExpectQuery(expectedUserQuery).WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{"id", "Username", "UserMailId", "UserRole", "Password"}).AddRow(
			1, "John Doe", "email@test.com", "Admin", "AdminPass1"))
	mock.ExpectQuery(expectedPrivilegeQuery).WillReturnRows(sqlmock.NewRows([]string{"PrivilegeType", "UserID"}).AddRow("All", 1))

	// Since admin user, TXN will be rolled back & actual Delete query will not be fired
	mock.ExpectBegin()
	mock.ExpectRollback()

	// Call the function
	usersDeleted := DeleteUserIfExists(db, model.User{Model: gorm.Model{ID: 1}})

	// Ensure all expectations were met if there is no admin user
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, -1, usersDeleted)
}
