package main

import (
	"fmt"
	"net/http"
	"todoapi/database"

	"github.com/gorilla/mux"
)

func main() {
	database.InitDB()

	r := mux.NewRouter()

	//punlic routers


	//protected routes




	fmt.Println("server is running on port 8081...")
	err := http.ListenAndServe(":8081", r)
	if err != nil{
		panic(err)
	}
}