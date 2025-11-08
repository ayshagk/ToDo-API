package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//creating a jwt token for a userid that has an expiry time limit, before needed again.
func GenerateJWT(userID uint) (string, error) { 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ 
		"user_id": userID,                                
		"exp":     time.Now().Add(time.Hour * 24).Unix(), 
	})

	secret := os.Getenv("JWT_SECRET")                      
	tokenString, err := token.SignedString([]byte(secret)) 
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//is the jwt token valid? parsing it means verifying its signature and decode claims so program can use it 
func VerifyJWT(tokenString string) (*jwt.Token, error) { 
	secret := os.Getenv("JWT_SECRET")                                      
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) { 
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { 
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"]) 
		}
		return []byte(secret), nil 
	})

	if err != nil {
		return nil, err 
	}

	return token, nil 

}

// protect routes and check every request has valid jwt token before it can happen/security
func AuthMiddleware(next http.Handler) http.Handler { 
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
		authHeader := r.Header.Get("Authorization") 
		if authHeader == "" {                       
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		
        //remove bearer from auth header to have token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ") 

		_, err := VerifyJWT(tokenString) 
		if err != nil {
			http.Error(w, "invalid token", http.StatusInternalServerError) 
			return
		}
		next.ServeHTTP(w, r) 

	})
}

// take userid from token to allow users with that id to make requests
//read auth header and make sure it exists
func GetUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization") 
	if authHeader == "" {                       
		return 0, fmt.Errorf("authorization header required")
	}

	//remove bearer before token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ") //remove bearer from token

	//verify token
	token, err := VerifyJWT(tokenString) 
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	//parse/verify the token and get the userid
	claims, ok := token.Claims.(jwt.MapClaims) 
	if !ok {
		return 0, fmt.Errorf("could not parse claims")
	}

	//if userid isnt a number/missing, an error will be returned 
	userID, ok := claims["user_id"].(float64) 
	if !ok {
		return 0, fmt.Errorf("user id not found in token")
	}
 
    //if the token has a ok userid then its returned as a uint
	return uint(userID), nil 

}
