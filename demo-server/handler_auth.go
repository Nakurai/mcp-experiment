package main

import (
	"fmt"
	"net/http"
)

func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	githubLoginUrl := makeGithubLoginUrl()
	http.Redirect(w, r, githubLoginUrl, http.StatusFound)
}
func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state != ENV["SERVER_KEY"] {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

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
