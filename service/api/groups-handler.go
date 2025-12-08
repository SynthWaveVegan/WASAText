package api

import (
	"encoding/json"
	"errors"
	"github.com/SynthWaveVegan/WASAText/service/api/reqcontext"
	"github.com/SynthWaveVegan/WASAText/service/structs"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func (rt *_router) ADDTOGROUP(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	groupId := ps.ByName("groupId")
	userId := ps.ByName("userId")

	if userId == "" || groupId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no groupId or userId retrieved")
		return
	}

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}

	var GroupName string
	err := json.NewDecoder(r.Body).Decode(&GroupName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
		return

	}
	defer r.Body.Close()

	UserId := structs.Identifier{
		Id: userId,
	}

	err = rt.db.AddToGroup(GroupName, UserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("User added to group successfully")
}

func (rt *_router) LEAVEGROUP(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")
	groupId := ps.ByName("groupId")

	if userId == "" || groupId == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no userId or groupId retrieved")
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

	GroupId := structs.Identifier{
		Id: groupId,
	}

	err := rt.db.LeaveGroup(GroupId, UserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Group deleted successfully")
}

func (rt *_router) SETGROUPNAME(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	groupId := ps.ByName("groupId")
	if groupId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var GroupName string
	err := json.NewDecoder(r.Body).Decode(&GroupName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
		return

	}
	defer r.Body.Close()

	GroupId := structs.Identifier{
		Id: groupId,
	}

	err = rt.db.SetGroupName("Update", GroupName, GroupId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Group name updated successfully")
}
