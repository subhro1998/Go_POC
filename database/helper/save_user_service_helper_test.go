package helper

import (
	"Go_Assignment/database/util"
	"Go_Assignment/dto"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveUserDetails_Success(t *testing.T) {
	db, mock, err := util.SetUpMockDB()
	assert.NoError(t, err)

	// Insert query
	mock.ExpectBegin()
	expectedSaveUserQuery := "INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`user_name`,`user_email`,`user_role`,`user_password`,`id`) VALUES (?,?,?,?,?,?,?,?)"
	expectedSavePrivilegeQuery := "INSERT INTO `user_privileges` (`created_at`,`updated_at`,`deleted_at`,`privilege_type`,`user_id`,`id`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`)"

	mock.ExpectExec(expectedSaveUserQuery).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sql.NullTime{Time: time.Now(), Valid: false}, "XYZ", "xyz@email.com", "Admin", "admin_pass", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(expectedSavePrivilegeQuery).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sql.NullTime{Time: time.Now(), Valid: false}, "All", 1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// User request
	userPrivilege := dto.UserPrivilegeResponse{
		UserId:        1,
		PrivilegeId:   2,
		PrivilegeType: "All",
	}
	userReq := dto.UserRequest{
		UserId:     1,
		Name:       "XYZ",
		Email:      "xyz@email.com",
		Password:   "admin_pass",
		Role:       "Admin",
		Privileges: []dto.UserPrivilegeResponse{userPrivilege},
	}

	userDataAbstract := make([]interface{}, 1)
	userDataAbstract[0] = userReq
	status, statusMessage := SaveUserDetails(db, userDataAbstract)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, status, http.StatusAccepted)
	assert.Equal(t, statusMessage, "Provided 1 User is saved in DB")
}
