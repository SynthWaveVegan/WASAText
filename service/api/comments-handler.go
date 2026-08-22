package api

import (
	"encoding/json"
	// "errors"
	"github.com/SynthWaveVegan/WASAText/service/api/reqcontext"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func (rt *_router) COMMENTMESSAGE(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")
	messageId := ps.ByName("messageId")

	if userId == "" || messageId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no userId or messageId retrieved")
		return
	}

	MessageId := structs.Identifier{
		Id: messageId,
	}
	UserId := structs.Identifier{
		Id: userId,
	}

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}

	type CommentRequest struct {
		CommentBody string `json:"CommentBody"`
	}

	var commentReq CommentRequest

	err := json.NewDecoder(r.Body).Decode(&commentReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
		return

	}
	defer r.Body.Close()

	log.Println("comment body: ", commentReq.CommentBody)

	if commentReq.CommentBody == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("comment body is empty")
		return
	}

	NewComment, err := rt.db.CommentMessage(commentReq.CommentBody, MessageId, UserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(NewComment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
	}

	log.Println("Message commented successfully")

}

func (rt *_router) UNCOMMENTMESSAGE(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")
	messageId := ps.ByName("messageId")
	commentId := ps.ByName("commentId")

	if userId == "" || messageId == "" || commentId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no userId or messageId retrieved")
		return
	}

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}
	CommentId := structs.Identifier{
		Id: commentId,
	}

	err := rt.db.UncommentMessage(CommentId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Comment deleted successfully")
}
