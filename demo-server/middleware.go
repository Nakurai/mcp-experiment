package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/nakurai/mcp-experiment/demo-server/internal/models"
	"github.com/nakurai/mcp-experiment/demo-server/internal/utils"
	"gorm.io/gorm"
)

func withDBSession(next http.Handler, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := db.WithContext(r.Context()) // creates a new session
		ctx := context.WithValue(r.Context(), utils.DbSessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := utils.GetDb(r.Context())
		if err != nil {
			http.Error(w, "No db provided", http.StatusInternalServerError)
			return
		}
		authorization, ok := r.Header["Authorization"]
		if !ok || len(authorization) == 0 {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}
		jwtBearer := strings.Split(authorization[0], " ")
		if len(jwtBearer) < 2 {
			http.Error(w, "Ill formatted token", http.StatusUnauthorized)
			return
		}
		jwt := jwtBearer[1]

		authToken, err := models.GetToken(db, jwt)
		if err != nil {
			http.Error(w, "provided token does not exist", http.StatusUnauthorized)
			return
		}

		err = utils.CheckJwt(authToken.Token, authToken.UserNumber)
		if err != nil {
			log.Default().Printf("in withAuth provided token is not valid: %v\n", err)
			http.Error(w, "provided token is not valid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", authToken.UserNumber)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
