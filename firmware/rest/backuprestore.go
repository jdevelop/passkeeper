package rest

import (
	"io"
	"net/http"
)

func (r *RESTServer) backupCredentials(w http.ResponseWriter, req *http.Request) {
	reader, err := r.credStorage.BackupStorage()
	if err != nil {
		errorResp(w, "Can't read credentials", &err)
	}
	w.Header().Set("Content-Type", "application/json")
	io.Copy(jsonHeaders(corsHeaders(w)), reader)
}

func (r *RESTServer) restoreCredentials(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseMultipartForm(10000); err != nil {
		errorResp(w, "can't parse multipart request", &err)
		return
	}
	file, _, err := req.FormFile("file")
	if err != nil {
		errorResp(w, "no file content", &err)
		return
	}
	defer file.Close()
	if err := r.credStorage.RestoreStorage(file); err != nil {
		errorResp(w, "can't restore storage", &err)
		return
	}
	corsHeaders(jsonHeaders(w)).Write(restored)
}
