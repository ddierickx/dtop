package main

import (
	"github.com/nu7hatch/gouuid"
)

type Authenticator struct {
	cfg      *DTopConfiguration
	sessions map[string]DTopUser
}

func NewAuthenticator(cfg *DTopConfiguration) *Authenticator {
	auth := new(Authenticator)
	auth.cfg = cfg
	auth.sessions = make(map[string]DTopUser)
	return auth
}

func (auth *Authenticator) Login(username string, password string) (bool, string) {
	hashed := auth.hashPassword(password)

	for _, user := range auth.cfg.Users {
		if user.Username == username && user.Password == hashed {
			token := auth.generateToken()
			auth.sessions[token] = user
			return true, token
		}
	}

	return false, ""
}

func (auth *Authenticator) Logout(token string) {
	delete(auth.sessions, token)
}

func (auth *Authenticator) IsAuthenticated(token string) bool {
	_, found := auth.sessions[token]
	return found
}

func (auth *Authenticator) generateToken() string {
	token, _ := uuid.NewV4()
	return token.String()
}

func (auth *Authenticator) hashPassword(password string) string {
	return password
}
