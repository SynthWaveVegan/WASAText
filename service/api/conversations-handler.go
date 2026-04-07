package api

import (
	"encoding/json"
	//"errors"
	"github.com/SynthWaveVegan/WASAText/service/api/reqcontext"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"fmt"
)

func (rt *_router) GETCONVERSATION(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationId := ps.ByName("conversationId")
	userId := ps.ByName("userId")

	fmt.Println("Conversation ID: ", conversationId)  // Log del conversationId
	fmt.Println("User ID: ", userId)  // Log del userId

	if conversationId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no conversationid retrieved")
		return
	}
	
	ConversationId := structs.Identifier{
		Id: conversationId,
	}

	ThisConversation, err := rt.db.GetConversation(ConversationId, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("1 something went wrong: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ThisConversation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("2 something went wrong: ", err)
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Conversation retrieved successfully")
}

func (rt *_router) CREATECONVERSATION(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	
	userId := ps.ByName("userId")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no userId retrieved")
		return
	}

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}
	
	UserId := structs.Identifier{
		Id: userId,
	}

	type RequestBody struct {
    	Name string `json:"Name"`
	}
	
	var username RequestBody
	
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("1 something went wrong: ", err)
		return

	}
	defer r.Body.Close()

	NewConversation, err := rt.db.CreateConversation(UserId, username.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("2 something went wrong: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(NewConversation)
	if err != nil {
    	ctx.Logger.Error("3 something went wrong: ", err)
    	return
	}

	log.Println("Conversation created successfully")

}

func (rt *_router) GETMYCONVERSATIONS(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no userId retrieved")
		return
	}

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}

	UserId := structs.Identifier{
		Id: userId,
	}

	MyConversations, err := rt.db.GetMyConversations(UserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(MyConversations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Conversations retrieved successfully")

}
