package handlers

import (
	"log"
	"net/http"

	"github.com/nakurai/mcp-experiment/demo-server/internal/models"
	"github.com/nakurai/mcp-experiment/demo-server/internal/statetoken"
	"github.com/nakurai/mcp-experiment/demo-server/internal/utils"
)

type GithubLoginUrlRes struct{
	Url string `json:"url"`
	Ok bool `json:"ok"`
	Message string `json:"message"`
}

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	githubLoginUrl, err := makeGithubLoginUrl()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := GithubLoginUrlRes{Message:"Internal Error"}
		utils.MakeRes(w, res)
		return
	}

	res := GithubLoginUrlRes{Url:githubLoginUrl, Ok: true}
	utils.MakeRes(w, res)
	// http.Redirect(w, r, githubLoginUrl, http.StatusFound)
}

type GithubCallbackRes struct{
	Token string `json:"token"`
	Ok bool `json:"ok"`
	Message string `json:"message"`
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {

	// at this point, the user authorized our app to use their email and we get a code from github that we can
	// exchange for an access token
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	err := statetoken.VerifyToken(state)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		res := GithubCallbackRes{Message:"Invalid Token - Cross Request attack likely - All is safe"}
		utils.MakeRes(w, res)
		// http.Error(w, "Invalid Token - Cross Request attack likely - All is safe", http.StatusUnauthorized)
		return
	}

	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		res := GithubCallbackRes{Message:"Invalid request"}
		utils.MakeRes(w, res)
		// http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// use the code provided by github to get
	token, err := getGithubAccessToken(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := GithubCallbackRes{Message:"github oauth failed "+err.Error()}
		utils.MakeRes(w, res)
		// http.Error(w, "github oauth failed "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use token to call GitHub API
	userEmail, err := getGithubUserEmail(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := GithubCallbackRes{Message:"github user info failed "+err.Error()}
		utils.MakeRes(w, res)
		// http.Error(w, "github user info failed "+err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := utils.GetDb(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := GithubCallbackRes{Message:"No db session provided"}
		utils.MakeRes(w, res)
		// http.Error(w, "No db session provided", http.StatusInternalServerError)
		return
	}

	// checkiing if the user exists or not
	user, err := models.GetUserByEmail(db, userEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := GithubCallbackRes{Message:"Error fetching our user"}
		utils.MakeRes(w, res)
		// http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	// here the user is new so we need to create them
	if user == nil {
		user, err = models.CreateUser(db, userEmail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := GithubCallbackRes{Message:"Error creating user"}
			utils.MakeRes(w, res)
			// http.Error(w, "Error creating the user", http.StatusInternalServerError)
			return
		}
	}

	// and now generating their jwt
	jwt, err := models.CreateToken(db, user.UserNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("Error: %v\n", err)
		res := GithubCallbackRes{Message:"Error creating the jwt"}
		utils.MakeRes(w, res)
		// http.Error(w, "Error creating the jwt", http.StatusInternalServerError)
		return
	}

	// Success!
	// fmt.Fprintf(w, "Authenticated as: %v\n%s", user, jwt)
	res := GithubCallbackRes{Ok:true, Token: jwt.Token}
	utils.MakeRes(w, res)

}
