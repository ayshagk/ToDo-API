package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Todos    []Todo `json:"todos" gorm:"foreignKey:UserID"`
}

type Todo struct {
	gorm.Model
	Title  string `json:"title"`
	Check  bool   `json:"check"`
	UserID uint   `json:"user_id"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
