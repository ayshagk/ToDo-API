package main

import (
	"fmt"
	"net/http"
	"todoapi/database"
	"todoapi/handlers"
	"todoapi/middleware"

	"github.com/gorilla/mux"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()

	//punlic routers
	r.HandleFunc("/register", handlers.RegisterUser)
	r.HandleFunc("/login", handlers.Login)

	//protected routes
	protectRoutes := r.PathPrefix("/").Subrouter() //create subrouter for routes that need authentication
	protectRoutes.Use(middleware.AuthMiddleware)

	protectRoutes.HandleFunc("/createTodo", handlers.CreateToDo)
	protectRoutes.HandleFunc("/getAllTodos", handlers.GetToDos)
	protectRoutes.HandleFunc("/updateTodo", handlers.UpdateToDo)
	protectRoutes.HandleFunc("/deleteTodo", handlers.DeleteToDo)

	fmt.Println("server is running on port 8081...")
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
