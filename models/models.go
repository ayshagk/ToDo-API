package models

import "gorm.io/gorm"

//created structs for all info 
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Todos    []Todo `json:"todos" gorm:"foreignKey:UserID"` //slice of todo structs inside user struct, one to many relationship as ONE USER CAN HAVE MANY TODOS.
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
