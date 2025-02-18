package helper

import (
	"Go_Assignment/database/util"
	"Go_Assignment/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllUsersSuccess(t *testing.T) {
	db, mock, err := util.SetUpMockDB()
	assert.NoError(t, err)

	// Define the expected queries
	expectedUserQuery := "SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL"
	expectedPrivilegeQuery := "SELECT * FROM `user_privileges` WHERE `user_privileges`.`user_id` = ? AND `user_privileges`.`deleted_at` IS NULL"

	// Set up the expected queries in the mock
	mock.ExpectQuery(expectedUserQuery).WillReturnRows(
		sqlmock.NewRows([]string{"id", "Username", "UserMailId", "UserRole", "Password"}).AddRow(
			1, "John Doe", "email@test.com", "Admin", "AdminPass1"))

	mock.ExpectQuery(expectedPrivilegeQuery).WillReturnRows(sqlmock.NewRows([]string{"PrivilegeType", "UserID"}).AddRow("All", 1))

	// Call the function
	userResponses := FetchAllUsers(db)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userResponses)
	assert.Equal(t, 1, len(userResponses))
	assert.Equal(t, "John Doe", userResponses[0].Name)
	assert.Equal(t, "email@test.com", userResponses[0].Email)
	assert.Equal(t, "Admin", userResponses[0].Role)

	assert.NotNil(t, userResponses[0].Privileges)
	assert.Equal(t, 1, len(userResponses[0].Privileges))
	assert.Equal(t, "All", userResponses[0].Privileges[0].PrivilegeType)
}

func TestFetchAllUsersNoUserFound(t *testing.T) {
	db, mock, err := util.SetUpMockDB()
	assert.NoError(t, err)

	// Define the expected queries
	expectedUserQuery := "SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL"

	// Set up the expected queries in the mock and return no response
	mock.ExpectQuery(expectedUserQuery).WillReturnRows(sqlmock.NewRows([]string{"id", "Username", "UserMailId", "UserRole", "Password"}))

	// Call the function
	userResponses := FetchAllUsers(db)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userResponses)
	assert.Equal(t, 0, len(userResponses))
}

func TestFetchUsersWithProvidedDetailsSuccess(t *testing.T) {
	db, mock, err := util.SetUpMockDB()
	assert.NoError(t, err)

	// Define the expected queries
	expectedUserQuery := "SELECT * FROM `users` WHERE `users`.`user_email` = ? AND `users`.`deleted_at` IS NULL"
	expectedPrivilegeQuery := "SELECT * FROM `user_privileges` WHERE `user_privileges`.`user_id` = ? AND `user_privileges`.`deleted_at` IS NULL"

	// Set up the expected queries in the mock
	mock.ExpectQuery(expectedUserQuery).WithArgs("email@test.com").WillReturnRows(
		sqlmock.NewRows([]string{"id", "Username", "UserMailId", "UserRole", "Password"}).AddRow(
			1, "John Doe", "email@test.com", "Admin", "AdminPass1"))
	mock.ExpectQuery(expectedPrivilegeQuery).WillReturnRows(sqlmock.NewRows([]string{"PrivilegeType", "UserID"}).AddRow("All", 1))

	// Call the function
	userResponses := FetchUsersWithProvidedDetails(db, model.User{UserMailId: "email@test.com"})

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userResponses)
	assert.Equal(t, 1, len(userResponses))
	assert.Equal(t, "John Doe", userResponses[0].Name)
	assert.Equal(t, "email@test.com", userResponses[0].Email)
	assert.Equal(t, "Admin", userResponses[0].Role)

	assert.NotNil(t, userResponses[0].Privileges)
	assert.Equal(t, 1, len(userResponses[0].Privileges))
	assert.Equal(t, "All", userResponses[0].Privileges[0].PrivilegeType)
}

func TestFetchUsersWithProvidedDetailsUserNotFound(t *testing.T) {
	db, mock, err := util.SetUpMockDB()
	assert.NoError(t, err)

	// Define the expected queries
	expectedUserQuery := "SELECT * FROM `users` WHERE `users`.`user_email` = ? AND `users`.`deleted_at` IS NULL"

	// Set up the expected queries in the mock
	mock.ExpectQuery(expectedUserQuery).WithArgs("email@test.com").WillReturnRows(
		sqlmock.NewRows([]string{"id", "Username", "UserMailId", "UserRole", "Password"}))

	// Call the function
	userResponses := FetchUsersWithProvidedDetails(db, model.User{UserMailId: "email@test.com"})

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.NotNil(t, userResponses)
	assert.Equal(t, 0, len(userResponses))
}
