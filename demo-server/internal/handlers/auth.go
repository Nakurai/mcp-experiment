package handlers

import (
	"fmt"
	"net/http"

	"github.com/nakurai/mcp-experiment/demo-server/internal/models"
	"github.com/nakurai/mcp-experiment/demo-server/internal/statetoken"
	"github.com/nakurai/mcp-experiment/demo-server/internal/utils"
)

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	githubLoginUrl, err := makeGithubLoginUrl()
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, githubLoginUrl, http.StatusFound)
}
func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {

	// at this point, the user authorized our app to use their email and we get a code from github that we can
	// exchange for an access token
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	err := statetoken.VerifyToken(state)
	if err != nil {
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

	db, err := utils.GetDb(r.Context())
	if err != nil {
		http.Error(w, "No db session provided", http.StatusInternalServerError)
		return
	}

	// checkiing if the user exists or not
	user, err := models.GetUserByEmail(db, userEmail)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	// here the user is new so we need to create them
	if user == nil {
		user, err = models.CreateUser(db, userEmail)
		if err != nil {
			http.Error(w, "Error creating the user", http.StatusInternalServerError)
			return
		}
	}

	// and now generating their jwt
	jwt, err := models.CreateToken(db, user.UserNumber)
	if err != nil {
		http.Error(w, "Error creating the jwt", http.StatusInternalServerError)
		return
	}

	// Success!
	fmt.Fprintf(w, "Authenticated as: %v\n%s", user, jwt)

}
