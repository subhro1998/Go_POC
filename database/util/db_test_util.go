package util

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetUpMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}

	dilact := mysql.New(
		mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		})

	gormDB, errGorm := gorm.Open(dilact, &gorm.Config{})
	if errGorm != nil {
		return nil, nil, err
	}

	return gormDB, sqlMock, nil
}
