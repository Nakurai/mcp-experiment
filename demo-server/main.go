package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"statetoken"
	"strings"
)

var ENV = map[string]string{}

func init() {
	// Load environement variables
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

func createApp(ctx context.Context) *http.ServeMux{
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handleGitHubLogin(w, r, ctx)
	} )
	mux.HandleFunc("/api/github/callback", func(w http.ResponseWriter, r *http.Request) {
		handleGitHubCallback(w, r, ctx)
	} )
	return mux
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start state token expiration checks
	statetoken.Init(ctx)

	// creating the app
	mux := createApp(ctx)
	app := &http.Server{Addr: ":8090", Handler: mux}
	
	// starting the server!
	go func() {
		fmt.Println("Starting the server on port 8090")
		if err := app.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()	

	// listening for interruption signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// cleaning up and shutting everything down nicely
	fmt.Println("shutting everything down!")
	cancel()
	app.Shutdown(ctx)
	fmt.Println("all done")
}
