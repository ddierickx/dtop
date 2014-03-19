package main

import (
	"testing"
)

func TestLogin(t *testing.T) {
	testUser := NewDTopUser("ho", "dor")
	users := []DTopUser{*testUser}
	cfg := NewDTopConfiguration("name", "description", users, "static", 12345)
	auth := NewAuthenticator(cfg)

	if ok, _ := auth.Login("mscott", "theboss"); ok {
		panic("Login succeeded but user does not exist.")
	}

	if ok, _ := auth.Login("ho", "dors"); ok {
		panic("Login succeeded but password was incorrect.")
	}

	var ok bool
	var token string

	if ok, token = auth.Login("ho", "dor"); !ok || token == "" {
		panic("Successful login failed or token was not provided.")
	}

	if authenticated := auth.IsAuthenticated(token); !authenticated {
		panic("User was not authenticated but should have been.")
	}

	auth.Logout(token)

	if authenticated := auth.IsAuthenticated(token); authenticated {
		panic("User logged out but still appears to be authenticated.")
	}
}
