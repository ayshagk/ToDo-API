package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint) (string, error) { //returns signed token string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //create new token with claims , and hs256 is encryption(secure signature)
		"user_id": userID,                                //claims are data in the token
		"exp":     time.Now().Add(time.Hour * 24).Unix(), //how long before new token needed
	})

	secret := os.Getenv("JWT_SECRET")                      //read secret key from env, extra secure.
	tokenString, err := token.SignedString([]byte(secret)) //sign token so it cannot be edited.
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) { //take jwt string and check if valid, return decoded token if valid.
	secret := os.Getenv("JWT_SECRET")                                      //check if token real or fake
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) { //open token check inside if signature match secret key (parse)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { //check if token signed with hs256. !ok check if its right type
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"]) //error if token isnt using right signing type
		}
		return []byte(secret), nil //return secret key back to parser so can verify token signature, if no match then invalid
	})

	if err != nil {
		return nil, err //return error if error while parsing like expired, invalid etc
	}

	return token, nil //jwt valid and tokn veriified

}

// protect routes and check every request has valid jwt token before it can happen/security
func AuthMiddleware(next http.Handler) http.Handler { //takes in next (actual route func).
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //return new http handler which runs before any protected endpoint
		authHeader := r.Header.Get("Authorization") //get bearer token
		if authHeader == "" {                       //if no token reject it so only logged in users can do actions.
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		// "Bearer eyjjbcjbdcjbdkj.cbkjBCJKBCKJEBFF.6G7347ED87GHE8DFU" example

		tokenString := strings.TrimPrefix(authHeader, "Bearer ") //remove bearer word from start of token

		_, err := VerifyJWT(tokenString) //use verifyjwt func to check if token valid
		if err != nil {
			http.Error(w, "invalid token", http.StatusInternalServerError) //stop request if invalid
			return
		}
		next.ServeHTTP(w, r) //if token valid then allow next handler to happen

	})
}

// take userid from token to connect posts/comments to user
func GetUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization") //get bearer token
	if authHeader == "" {                       //if none then error
		return 0, fmt.Errorf("authorization header required")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ") //remove bearer from token

	token, err := VerifyJWT(tokenString) //check if token valid like before
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims) //converting the claims data into mapclaims so we can access
	if !ok {
		return 0, fmt.Errorf("could not parse claims")
	}

	userID, ok := claims["user_id"].(float64) //get userid from claims so we can then..
	if !ok {
		return 0, fmt.Errorf("user id not found in token")
	}

	return uint(userID), nil //uint back to handler so route knows which user made request to add into column

}
