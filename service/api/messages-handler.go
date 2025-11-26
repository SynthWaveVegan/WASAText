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

	var MessageBody string

	err := json.NewDecoder(r.Body).Decode(&MessageBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("something went wrong: ", err)
		return

	}
	defer r.Body.Close()
}