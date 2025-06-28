package handlers

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/nakurai/mcp-experiment/demo-server/internal/models"
	"github.com/nakurai/mcp-experiment/demo-server/internal/utils"
)

type NewMessageBody struct {
	Content      string `json:"content"`
	ContactEmail string `json:"contactEmail"`
	Tag          string `json:"tag"`
}

type NewMessageRes struct {
	Ok      bool
	Message string
}

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		res := NewMessageRes{Message: "only POST request are allowed"}
		utils.MakeRes(w, res)
		return
	}
	newMessageInfo := NewMessageBody{}
	err := json.NewDecoder(r.Body).Decode(&newMessageInfo)
	if err != nil {
		log.Default().Printf("error unmarshalling message info: %#v", err)
		w.WriteHeader(http.StatusBadRequest)
		res := NewMessageRes{Message: "Your request is ill formatted"}
		utils.MakeRes(w, res)
		return
	}

	userNumber := r.Context().Value("user").(string)
	newMessage := &models.Message{
		UserNumber:   userNumber,
		ContactEmail: strings.TrimSpace(html.EscapeString(newMessageInfo.ContactEmail)),
		Content:      strings.TrimSpace(html.EscapeString(newMessageInfo.Content)),
		Tag:          strings.TrimSpace(html.EscapeString(newMessageInfo.Tag)),
	}

	db, err := utils.GetDb(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res := NewMessageRes{Message: "In HandleMessage accessing db session"}
		utils.MakeRes(w, res)
		return
	}

	_, err = models.CreateMessage(db, newMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Default().Printf("in HandleMessage: %v", err)
		res := NewMessageRes{Message: "Error while saving the message"}
		utils.MakeRes(w, res)
		return
	}

	res := NewMessageRes{Ok: true}
	utils.MakeRes(w, res)
}
