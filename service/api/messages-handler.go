package api

import (
	"github.com/SynthWaveVegan/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"encoding/json"
	"errors"
	"log"
)

func (rt * router) SENDMESSAGE(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	
	userId := ps.ByName("userId")
	messageId := ps.ByName("messageId")

	if userId == "" || messageId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}

	MediaType := r.Header.Get("MediaType")

	if MediaType == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	} else if MediaType != "Text" && MediaType != "Photo"{
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return

	var MessageBody string

	err := json.NewDecoder(r.Body).Decode(&MessageBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
		return

	}
	defer r.Body.Close()

	NewMessage, err := rt.db.sendMessage(MessageBody, userId, MediaType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Message sent successfully")
	
}