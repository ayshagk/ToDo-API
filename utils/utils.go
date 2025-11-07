package utils

import "golang.org/x/crypto/bcrypt"

//create hashed pass using generatefrompassword
func HashPassword(password string) (string, error) {  //plain password string as input and return hashed pass string and error if wrong
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //convert pass string into byte as bycrpyt requires this
	if err != nil {                                                   //defaultcost (higher cost-more secure)            
		return "", err
	}

	return string(HashPassword), nil    //convert hash pass from byte slice to string and return it
}


func ComparePassword(plainPass, HashedPass string) error { //take in plain pass and hashed to compare
	err := bcrypt.CompareHashAndPassword([]byte(HashedPass), []byte(plainPass)) //comparefunc to compare passwords and convert to byte slices.
	if err != nil {
		return err
	}
	return nil                   //nil if plain matches hashed and error if dont
}