package api

import (
	"encoding/json"
	//"errors"
	"github.com/SynthWaveVegan/WASAText/service/api/reqcontext"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type MessageRequest struct {
  MessageBody string `json:"messageBody"`
}

func (rt *_router) SENDMESSAGE(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
  
  userId := ps.ByName("userId")
  conversationId := ps.ByName("conversationId")

  if userId == "" || conversationId == "" {
    w.WriteHeader(http.StatusBadRequest)
    ctx.Logger.Error("no userId or conversationId retrieved")
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
    ctx.Logger.Error("no mediatype retrieved")
    return
  } else if MediaType != "Text" && MediaType != "Photo" {
    w.WriteHeader(http.StatusBadRequest)
    ctx.Logger.Error("mediatype not valid for use")
    return
  }

  var requestBody MessageRequest

  err := json.NewDecoder(r.Body).Decode(&requestBody)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    ctx.Logger.Error("something went wrong: ", err)
    return
  }
  defer r.Body.Close()

  MessageBody := requestBody.MessageBody

  UserId := structs.Identifier{
    Id: userId,
  }

  ConversationId := structs.Identifier{
    Id: conversationId,
  }

  NewMessage, err := rt.db.SendMessage(MessageBody, UserId, MediaType, ConversationId)
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

func (rt *_router) FORWARDMESSAGE(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")
	messageId := ps.ByName("messageId")
	conversationId := ps.ByName("conversationId")

	if userId == "" || messageId == "" || conversationId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no userId or messageId or conversationId retrieved")
		return
	}

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}

	var TargetConversationId structs.Identifier

	err := json.NewDecoder(r.Body).Decode(&TargetConversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
		return

	}
	defer r.Body.Close()

	UserId := structs.Identifier{
		Id: userId,
	}

	MessageId := structs.Identifier{
		Id: messageId,
	}

	ForwardedMessage, err := rt.db.ForwardMessage(MessageId, UserId, TargetConversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ForwardedMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Message forwarded successfully")

}

func (rt *_router) DELETEMESSAGE(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")
	messageId := ps.ByName("messageId")
	conversationId := ps.ByName("conversationId")

	if userId == "" || messageId == "" || conversationId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no userId or messageId or conversationid retrieved")
		return
	}

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}

	MessageId := structs.Identifier{
		Id: messageId,
	}

	err := rt.db.DeleteMessage(MessageId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Message deleted successfully")
}

func (rt *_router) MARKMESSAGEREAD(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")
	conversationId := ps.ByName("conversationId")

	if userId == "" || conversationId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no userId or messageId or conversationid retrieved")
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

	ConversationId := structs.Identifier{
		Id: conversationId,
	}

	err := rt.db.MarkMessageRead(UserId, ConversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Message deleted successfully")
}
