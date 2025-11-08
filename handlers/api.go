package handlers

import (
	"encoding/json"
	"net/http"
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
 if err != nil{
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
if err != nil{
	http.Error(w, "invalid password", http.StatusBadRequest)
	return 
}

//generate a jwt token
token, err:=middleware.GenerateJWT(user.ID)
if err != nil{
	http.Error(w, "failed to generate token", http.StatusInternalServerError)
return
}

w.WriteHeader(200)
json.NewEncoder(w).Encode(token)
}
/*
func CreateToDo(){

}

func GetToDos(){

}

func UpdateToDo(){


}

func DeleteToDo(){

}
*/