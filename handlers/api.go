package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todoapi/database"
	"todoapi/middleware"
	"todoapi/models"
	"todoapi/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req *models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	err = database.Db.Where("email = ?", req.Email).First(&user).Error
	if err == nil {
		http.Error(w, "sorry user already exists", http.StatusBadRequest)
		return
	}

	HashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "unable to hash the password", http.StatusBadRequest)
		return
	}

	req.Password = HashPassword

	err = database.Db.Create(&req).Error
	if err != nil {
		http.Error(w, "unable to create user", http.StatusBadRequest)
		return
	}

	// send a response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user registered successfully!"))

}

func Login(w http.ResponseWriter, r *http.Request) {
	//decode the request from the request body which is the email and pass
	var login models.Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//check if the user exists when logging on
	var user models.User
	err = database.Db.Where("email = ?", login.Email).First(&user).Error
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	//check if password matches the one in our db
	err = utils.ComparePassword(login.Password, user.Password)
	if err != nil {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	//generate a jwt token
	token, err := middleware.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)
}

func CreateToDo(w http.ResponseWriter, r *http.Request) {
	var addTodo models.Todo
	err := json.NewDecoder(r.Body).Decode(&addTodo)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addTodo.UserID = userID

	err = database.Db.Create(&addTodo).Error
	if err != nil {
		http.Error(w, "failed to add toDo", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "toDo added successfully:)")

}

func GetToDos(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "access denied ", http.StatusUnauthorized)
		return
	}
	var allTodos []models.Todo
	err = database.Db.Where("user_id = ?", userID).Find(&allTodos).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(allTodos)

}

func UpdateToDo(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	todoIDStr := r.URL.Query().Get("id")
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	err = database.Db.First(&todo, "id = ? AND user_id = ?", todoID, userID).First(&todo).Error
	if err != nil {
		http.Error(w, "this todo does not belong to the user", http.StatusNotFound)
		return
	}

	var update models.Todo
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if update.Title != "" {
		todo.Title = update.Title
	}

	todo.Check = update.Check
	database.Db.Save(&todo)

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todo)

}

func DeleteToDo(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	todoIDStr := r.URL.Query().Get("id")
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	result := database.Db.Delete(&models.Todo{}, "id = ? AND user_id = ?", todoID, userID)

	if result.RowsAffected == 0 {
		http.Error(w, "todo does not belong to the user or not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(200)
	fmt.Fprintln(w, "todo deleted successfully")

}
