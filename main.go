package main

//live app deployed on render : https://todo-api-gxa2.onrender.com

import (
	"fmt"
	"net/http"
	"os"
	"todoapi/database"
	"todoapi/handlers"
	"todoapi/middleware"

	"github.com/gorilla/mux"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()
// add this so it can maybe work on render!
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "API is live!")
})

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

	/*fmt.Println("server is running on port 8081...")
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	} */
//took away our port and added something that would allow it to connect to render.
	port := os.Getenv("PORT")
if port == "" {
    port = "8081" // fallback for local testing
}

fmt.Println("Server is running on port " + port)
err := http.ListenAndServe(":"+port, r)
if err != nil {
    panic(err)
}
}
