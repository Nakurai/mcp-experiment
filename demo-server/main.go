package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

var ENV = map[string]string{}

func init() {
	content, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	content_str := string(content)
	envs := strings.Split(content_str, "\n")
	for _, env := range envs {
		parts := strings.Split(env, "=")
		ENV[parts[0]] = parts[1]
	}
}

func makeRes(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/api/login", handleGitHubLogin)
	http.HandleFunc("/api/github/callback", handleGitHubCallback)
	http.ListenAndServe(":8090", nil)
}
