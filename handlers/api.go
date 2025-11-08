package handlers

import (
	"encoding/json"
	"net/http"
	"todoapi/database"
	"todoapi/models"
	"todoapi/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request){
	var req *models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	err = database.Db.Where("email = ?", req.Email).First(&user).Error
	 if err != nil{
		http.Error(w, "sorry user already exists", http.StatusBadRequest)
		return 
	 }

	HashPassword, err := utils.HashPassword(req.Password)
	if err != nil{
		http.Error(w, "unable to hash the password", http.StatusBadRequest)
		return 
	}

	req.Password = HashPassword

	
	 
}