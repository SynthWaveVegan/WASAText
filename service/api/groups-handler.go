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

func (rt * router) ADDTOGROUP(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	err := rt.db.addToGroup(GroupId structs.Identifier, AddUserId structs.Identifier)
}

func (rt * router) CREATEGROUP(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	NewGroup, err := rt.db.createGroup(GroupName string, UserId structs.Identifier)
}

func (rt * router) LEAVEGROUP(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	err := rt.db.leaveGroup(GroupId structs.Identifier, UserId structs.Identifier)
}

func (rt * router) SETGROUPNAME(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	err := rt.db.setGroupName("Update", GroupName string, groupId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Group name updated successfully")
}