package models

import (
	"golang.org/x/crypto/bcrypt"	
)

type(
	User struct{
		ID					int 			`json:"id,omitempty"`
		Username		string		`json:"username,omitempty"`
		Email 			string		`json:"email,omitempty"`
		Password 		string		`json:"password,omitempty"`
		Token 			string		`json:"jwt,omitempty"`
		Followers 	int				`json:"followers"`
	}
)

func (u *User) HashPassword() {
	bpwd := []byte(u.Password)
	password, err := bcrypt.GenerateFromPassword(bpwd, 10)
	if err == nil {
		u.Password = string(password)
	}
}