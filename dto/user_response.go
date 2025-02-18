package dto

type UserResponse struct {
	UserId      uint                    `json:"user_id"`
	Name        string                  `json:"user_name"`
	Email       string                  `json:"user_email"`
	Role        string                  `json:"user_role"`
	City        string                  `json:"city"`
	Zipcode     string                  `json:"zipcode"`
	DateofBirth string                  `json:"date_of_birth"`
	Password    string                  `json:"-"`
	Privileges  []UserPrivilegeResponse `json:"user_privileges"`
}

type UserRequest struct {
	UserId      uint                    `json:"user_id"`
	Name        string                  `json:"user_name"`
	Email       string                  `json:"user_email"`
	Role        string                  `json:"user_role"`
	City        string                  `json:"city"`
	Zipcode     string                  `json:"zipcode"`
	DateofBirth string                  `json:"date_of_birth"`
	Password    string                  `json:"password"`
	Privileges  []UserPrivilegeResponse `json:"user_privileges"`
}

type UserPrivilegeResponse struct {
	PrivilegeId   uint   `json:"privilege_id"`
	PrivilegeType string `json:"privilege_name"`
	UserId        uint   `json:"user_id"`
}
