package util

import (
	"fmt"
	"log"

	"Go_Assignment/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const dns = "root:root@tcp(127.0.0.1:3306)/USER_DB?charset=utf8mb4&parseTime=True&loc=Local"

func CreateConnection(autoCreateTables bool, enableGORMLogging bool) *gorm.DB {
	var config *gorm.Config
	if enableGORMLogging {
		config = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		config = &gorm.Config{}
	}

	db, errConn := gorm.Open(mysql.Open(dns), config)

	if errConn != nil {
		log.Fatal("Issue in opening new DB connection")
		panic(errConn)
	}

	if autoCreateTables {
		fmt.Println("Auto creating tables")
		db.AutoMigrate(&model.User{}, &model.UserPrivilege{})
	}

	return db
}

func CloseDBConnection(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal("Not able to do db.DB() to get open sql Connections")
		panic(db)
	}

	sqlDb.Close()
}
