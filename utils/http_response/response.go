package httpresponse

import (
	"encoding/json"
	"net/http"
)

func WriteMessage(w http.ResponseWriter, status int, msg string) {
	var j struct {
		Msg string `json:"message"`
	}

	j.Msg = msg

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func WriteResponse(w http.ResponseWriter, status int, data any) {
	var j struct {
		Data any `json:"data"`
	}

	j.Data = data

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

}

func WriteError(w http.ResponseWriter, status int, err error) {
	http.Error(w, err.Error(), status)
}
