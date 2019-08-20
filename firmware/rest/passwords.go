package rest

import (
	"bytes"
	"encoding/json"
	"github.com/jdevelop/passkeeper/firmware/storage"
	"io/ioutil"
	"net/http"
	"strings"
)

func (r *RESTServer) generatePasswords(w http.ResponseWriter, req *http.Request) {
	pwds, err := r.passwordGen.GeneratePassword(5)
	if err != nil {
		errorResp(w, "can't generate passwords", &err)
		return
	}
	respJson, err := json.Marshal(pwds)
	if err != nil {
		errorResp(w, "can't restore stora", &err)
		return
	}
	corsHeaders(jsonHeaders(w)).Write(respJson)
}

type CardPasswordUpdate struct {
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

const blockSize = 32

func (r *RESTServer) setPassword(w http.ResponseWriter, req *http.Request) {
	var passUpdate CardPasswordUpdate
	if err := json.NewDecoder(req.Body).Decode(&passUpdate); err != nil {
		errorResp(w, "can't unmarshal password update", &err)
		return
	}
	defer req.Body.Close()
	if strings.TrimSpace(passUpdate.Password) == "" || passUpdate.Password != passUpdate.Confirm {
		errorResp(w, "password can't be empty or mismatch", nil)
		return
	}

	block := storage.BuildKey([]byte(passUpdate.Password))

	rdr, err := r.credStorage.BackupStorage()
	if err != nil {
		errorResp(w, "can't read the storage content, aborting", &err)
		return
	}
	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		errorResp(w, "can't read the storage stream, aborting", &err)
		return
	}
	if err := r.cardPasswordChange(block); err != nil {
		errorResp(w, "can't update password, please try again", nil)
		return
	}
	if err := r.credStorage.UpdateKey(block); err != nil {
		errorResp(w, "can't set crypto storage password, aborting", nil)
		return
	}
	if err := r.credStorage.RestoreStorage(bytes.NewReader(data)); err != nil {
		errorResp(w, "can't restore password, aborting", nil)
		return
	}

	corsHeaders(jsonHeaders(w)).Write(saved)
}
