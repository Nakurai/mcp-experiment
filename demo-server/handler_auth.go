package main

import (
	"context"
	"fmt"
	"net/http"
	"statetoken"
)

func handleGitHubLogin(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	githubLoginUrl, err := makeGithubLoginUrl()
	if err != nil{
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, githubLoginUrl, http.StatusFound)
}
func handleGitHubCallback(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	// at this point, the user authorized our app to use their email and 
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	err := statetoken.VerifyToken(state)
	if err != nil{
		http.Error(w, "Invalid Token - Cross Request attack likely - All is safe", http.StatusUnauthorized)
		return
	}

	if code == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// use the code provided by github to get 
	token, err := getGithubAccessToken(code)
	if err != nil {
		http.Error(w, "github oauth failed "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use token to call GitHub API
	userEmail, err := getGithubUserEmail(token)
	if err != nil {
		http.Error(w, "github user info failed "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Success!
	fmt.Fprintf(w, "Authenticated as: %s", userEmail)
}
