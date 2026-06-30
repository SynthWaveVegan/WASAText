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

func (rt *_router) SETMYPHOTO(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type PhotoRequest struct {
    	Photo string `json:"Photo"`
	}

	var req PhotoRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	PhotoPath := req.Photo

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}

	UserId := structs.Identifier{
		Id: userId,
	}

	err = rt.db.SetMyPhoto(UserId, PhotoPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Profile Photo updated successfully")
}

func (rt *_router) SETGROUPPHOTO(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	groupId := ps.ByName("groupId")

	if groupId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type PhotoRequest struct {
    	Photo string `json:"Photo"`
	}

	var req PhotoRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	PhotoPath := req.Photo

	authorization := r.Header.Get("Authorization")

	if groupId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}

	GroupId := structs.Identifier{
		Id: groupId,
	}

	err = rt.db.SetGroupPhoto(PhotoPath, GroupId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Group Photo updated successfully")
}
