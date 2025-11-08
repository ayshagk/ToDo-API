package utils

import "golang.org/x/crypto/bcrypt"

//create hashed pass using generatefrompassword
func HashPassword(password string) (string, error) {  
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //convert pass string into byte as bycrpyt requires this
	if err != nil {                                                              
		return "", err
	}

	return string(HashPassword), nil    //convert hash pass from byte slice to string and return it
}

//comparefunc to compare passwords and convert to byte slices.
func ComparePassword(plainPass, HashedPass string) error { 
	err := bcrypt.CompareHashAndPassword([]byte(HashedPass), []byte(plainPass)) 
	if err != nil {
		return err
	}
	return nil                   //nil if plain matches hashed and error if dont
}