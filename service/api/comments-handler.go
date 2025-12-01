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

func (rt * router) COMMENTMESSAGE(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	
	userId := ps.ByName("userId")
	messageId := ps.ByName("commentId")

	if userId == "" || commentId == "" {
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

	var CommentBody string

	err := json.NewDecoder(r.Body).Decode(&CommentBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
		return

	}
	defer r.Body.Close()

	Newcomment, err := rt.db.commentMessage(MessageBody, userId, MediaType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(NewMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Message sent successfully")
	
}