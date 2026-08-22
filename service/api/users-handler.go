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

func (rt *_router) GETUSER(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ThisUserId := structs.Identifier{
		Id: userId,
	}

	ThisUser, err := rt.db.GetUser(ThisUserId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("something went wrong: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ThisUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
	}

	w.WriteHeader(http.StatusOK)
	log.Println("User retrieved successfully")
}

func (rt *_router) SETMYUSERNAME(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type RequestBody struct {
		Name string `json:"Name"`
	}

	var username RequestBody

	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("5 something went wrong: ", err)
		return

	}
	defer r.Body.Close()

	authorization := r.Header.Get("Authorization")

	if userId != authorization {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("user is not allowed")
		return
	}

	err = rt.db.SetMyUsername("Update", username.Name, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("6 something went wrong: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Username updated successfully")

}
