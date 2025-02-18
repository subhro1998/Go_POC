package util

import (
	"Go_Assignment/dto"
	"Go_Assignment/model"
)

func ToResponseDto(user model.User) dto.UserResponse {
	userPrivilegesResponse := []dto.UserPrivilegeResponse{}
	for _, privilege := range user.Privileges {
		userPrivilege := dto.UserPrivilegeResponse{
			PrivilegeId:   privilege.ID,
			PrivilegeType: privilege.PrivilegeType,
			UserId:        privilege.UserID,
		}
		userPrivilegesResponse = append(userPrivilegesResponse, userPrivilege)
	}

	return dto.UserResponse{
		UserId:      user.ID,
		Name:        user.Username,
		Email:       user.UserMailId,
		Role:        user.UserRole,
		Privileges:  userPrivilegesResponse,
		Password:    user.Password,
		City:        "",
		Zipcode:     "",
		DateofBirth: "",
	}
}

func ToDomainDto(userRequest dto.UserRequest) model.User {
	userPrivilegesRequest := []model.UserPrivilege{}
	for _, privilege := range userRequest.Privileges {
		userDomainPrivilege := model.UserPrivilege{}
		if privilege.PrivilegeId > 0 { // Will take as existing Privilege record else will insert a new one
			userDomainPrivilege.ID = privilege.PrivilegeId
		}

		userDomainPrivilege.PrivilegeType = privilege.PrivilegeType
		userDomainPrivilege.UserID = privilege.UserId

		userPrivilegesRequest = append(userPrivilegesRequest, userDomainPrivilege)
	}

	domainUser := model.User{}
	if userRequest.UserId > 0 { // Will take as existing User record else will insert a new User
		domainUser.ID = userRequest.UserId
	}
	domainUser.Username = userRequest.Name
	domainUser.UserMailId = userRequest.Email
	domainUser.UserRole = userRequest.Role
	domainUser.Privileges = userPrivilegesRequest
	domainUser.Password = userRequest.Password

	return domainUser
}
