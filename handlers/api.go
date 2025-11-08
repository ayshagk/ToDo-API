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
//register user, decode request body into user struct first so can be read
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req *models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//does user already exist? check with email
	var user models.User
	err = database.Db.Where("email = ?", req.Email).First(&user).Error
	if err == nil {
		http.Error(w, "sorry user already exists", http.StatusBadRequest)
		return
	}

	//call on hashpassword to secure pass/hash
	HashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "unable to hash the password", http.StatusBadRequest)
		return
	}

	//make pass from request body into hashed pass
	req.Password = HashPassword

	//add user into db 
	err = database.Db.Create(&req).Error
	if err != nil {
		http.Error(w, "unable to create user", http.StatusBadRequest)
		return
	}

	// send a response to client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user registered successfully!"))

}
//user login to authenticate (recieve jwt token)
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

	//check if password matches the one in our db so verify
	err = utils.ComparePassword(login.Password, user.Password)
	if err != nil {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	//generate a jwt token- authentication 
	token, err := middleware.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
     
	//response and return token by encode (send back)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)
}

//add todo for user who must be authenticated 
func CreateToDo(w http.ResponseWriter, r *http.Request) {
	//decode todo details mentioned in the struct from request body.
	var addTodo models.Todo
	err := json.NewDecoder(r.Body).Decode(&addTodo)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	//user id from jwt token to be able to add
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//make todo's userid into logged in userid
	addTodo.UserID = userID

    //save to db table
	err = database.Db.Create(&addTodo).Error
	if err != nil {
		http.Error(w, "failed to add toDo", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "toDo added successfully:)")

}

//show all todos to client, only if authenticated
func GetToDos(w http.ResponseWriter, r *http.Request) {
	//user id from token
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "access denied ", http.StatusUnauthorized)
		return
	}

	//find all the todos that belong to this user. (where=which user)(find=look for)
	var allTodos []models.Todo
	err = database.Db.Where("user_id = ?", userID).Find(&allTodos).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//return to client as json so it can be read.
	json.NewEncoder(w).Encode(allTodos)

}

//update a todo if user authenticated
func UpdateToDo(w http.ResponseWriter, r *http.Request) {
	//token to get userid
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//get specific todoID, convert from string to int so it can be read (Atoi)
	todoIDStr := r.URL.Query().Get("id")
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	//make sure to get todo which belongs to this userid
	var todo models.Todo
	err = database.Db.First(&todo, "id = ? AND user_id = ?", todoID, userID).First(&todo).Error
	if err != nil {
		http.Error(w, "this todo does not belong to the user", http.StatusNotFound)
		return
	}

	//decode the updated info from request body that client wrote.
	var update models.Todo
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	//if client didnt write anything in title, then leave the old title info, instead of making it empty now.
	if update.Title != "" {
		todo.Title = update.Title
	}

	//update old check info with new info sent by the client.
	todo.Check = update.Check

	database.Db.Save(&todo) //save update into db

	//success message and encode to return 
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todo)

}

//delete specific todo if user is authenticated only
func DeleteToDo(w http.ResponseWriter, r *http.Request) {
	//user id from token
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//get the todoid from query parameter that client inputs in URL and convert to int so it can be read 
	todoIDStr := r.URL.Query().Get("id")
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	//delete the todo only if it belongs to the user with specific id
	result := database.Db.Delete(&models.Todo{}, "id = ? AND user_id = ?", todoID, userID)

	//if no rows are deleted then there is an error as the id didnt belong or was invalid
	if result.RowsAffected == 0 {
		http.Error(w, "todo does not belong to the user or not found", http.StatusNotFound)
		return
	}

	//success message and message to terminal hence Fprintln
	w.WriteHeader(200)
	fmt.Fprintln(w, "todo deleted successfully")

}


/*this file holds all the api endpoints for the Todo list application:
register, login, create, create, get all, update, and delete.
*/