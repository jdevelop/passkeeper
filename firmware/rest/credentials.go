package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jdevelop/passkeeper/firmware"
)

func (r *RESTServer) listCredentials(w http.ResponseWriter, _ *http.Request) {
	services, err := r.credStorage.ListCredentials()
	if err != nil {
		errorResp(w, "Can't load seeds", &err)
		return
	}
	data, err := json.Marshal(services)
	if err != nil {
		errorResp(w, "Can't marshal seeds", &err)
		return
	}
	_, err = jsonHeaders(corsHeaders(w)).Write(data)
	if err != nil {
		errorResp(w, "Response failure", &err)
		return
	}
	return
}

func (r *RESTServer) loadCredentials(w http.ResponseWriter, req *http.Request) {
	data := mux.Vars(req)
	if v, ok := data["id"]; !ok {
		if !ok {
			errorResp(w, "No required parameter", nil)
			return
		}
	} else {
		seed, err := r.credStorage.ReadCredentials(v)
		if err != nil {
			errorResp(w, "Can't find seed", &err)
			return
		}
		data, err := json.Marshal(&seed)
		if err != nil {
			errorResp(w, "Can't marshal seed", &err)
			return
		}
		jsonHeaders(corsHeaders(w)).Write(data)
	}
}

func (r *RESTServer) saveCredentials(w http.ResponseWriter, req *http.Request) {

	data, err := ioutil.ReadAll(req.Body)

	if err != nil {
		errorResp(w, "Can't read seed object", &err)
		return
	}

	var s firmware.Credentials

	if err = json.Unmarshal(data, &s); err != nil {
		errorResp(w, "Can't unmarshal seed object", &err)
		return
	}

	if s.Id == "" {
		s.Id = uuid.New().String()
	}

	if err = r.credStorage.WriteCredentials(s); err != nil {
		errorResp(w, "Can't save seed", &err)
		return
	}

	corsHeaders(jsonHeaders(w)).Write(saved)
	return
}

func (r *RESTServer) removeCredentials(w http.ResponseWriter, req *http.Request) {
	data := mux.Vars(req)
	if v, ok := data["id"]; !ok {
		errorResp(w, "No required parameter", nil)
		return
	} else {
		if err := r.credStorage.RemoveCredentials(v); err != nil {
			errorResp(w, "Can't remove seed", &err)
			return
		}
	}
	corsHeaders(jsonHeaders(w)).Write(removed)
}
