package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nakurai/mcp-experiment/demo-server/internal/handlers"
	"github.com/nakurai/mcp-experiment/demo-server/internal/models"
	"github.com/nakurai/mcp-experiment/demo-server/internal/statetoken"
)

var ENV = map[string]string{}
var db *gorm.DB

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


func createApp() *http.ServeMux{
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)
	mux.HandleFunc("/api/login", handlers.HandleGitHubLogin)
	mux.HandleFunc("/api/github/callback", handlers.HandleGitHubCallback)
	return mux
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	db, err = gorm.Open(sqlite.Open(ENV["DB_PATH"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.User{})

	// start state token expiration checks
	statetoken.Init(ctx)

	// creating the app
	mux := createApp()
	handler := withDBSession(mux, db)
	app := &http.Server{Addr: ":8090", Handler: handler}
	
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
