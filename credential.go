package main


const (
	noAuth byte = 0
	passAuth byte = 2
)



type AuthCredential struct {
	//map user -> password
	users map[string]string
}



// check if the client authentication is valid
func (ac *AuthCredential) Valid(user, pw string) bool {
	if len(ac.users) == 0 {
		//if no user exists no login is require
		return true
	}

	// if users exist check for passwors
	if pass, ok := ac.users[user]; ok && pass == pw {
		return true
	}

	// if login fail return false
	return false
}



// return a byte with the server auth method
func (ac *AuthCredential) GetServerMethod() byte {
	if len(ac.users) != 0 {
		return passAuth
	}

	//if no user are saved use the noAuth Method
	return noAuth
}