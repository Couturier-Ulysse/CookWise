package models

import "fmt"

type User struct {
	ID 		 int
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user *User) Create(Name string, Email string, Password string) error {
	user.Name = Name
	user.Email = Email
	user.Password = Password
	return nil
}

func (user *User) Get() (int, string, string, string) {
	return user.ID, user.Name, user.Email, user.Password
}

func (user *User) PrintInfo() {
	fmt.Println(&user.ID, &user.Name, &user.Email, &user.Password)
}