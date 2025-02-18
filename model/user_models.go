package model

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string          `gorm:"column:user_name"`
	UserMailId string          `gorm:"unique;column:user_email"`
	UserRole   string          `gorm:"column:user_role"`
	Password   string          `gorm:"not null;column:user_password"`
	Privileges []UserPrivilege `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserPrivilege struct {
	gorm.Model
	PrivilegeType string `gorm:"column:privilege_type;not null"`
	UserID        uint   `gorm:"column:user_id;not null"`
}

// BeforeDelete is a GORM hook that is called before a record is deleted
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if strings.EqualFold(u.UserRole, "admin") {
		return errors.New("admin user not allowed to be deleted")
	}
	return nil
}
