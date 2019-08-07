package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jdevelop/passkeeper/firmware"
	"github.com/jdevelop/passkeeper/firmware/storage"
)

type storageCombined interface {
	storage.CredentialsStorageList
	storage.CredentialsStorageRead
	storage.CredentialsStorageRemove
	storage.CredentialsStorageWrite
}

type RESTServer struct {
	credStorage storageCombined
}

func corsHeaders(w http.ResponseWriter) http.ResponseWriter {
	hdr := w.Header()
	hdr.Set("Access-Control-Allow-Origin", "*")
	hdr.Set("Access-Control-Allow-Methods", "OPTIONS,GET,PUT,DELETE,POST")
	hdr.Set("Access-Control-Allow-Headers", "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
	return w
}

func jsonHeaders(w http.ResponseWriter) http.ResponseWriter {
	hdr := w.Header()
	hdr.Set("Content-Type", "application/json")
	return w
}

func textHeaders(w http.ResponseWriter) http.ResponseWriter {
	hdr := w.Header()
	hdr.Set("Content-Type", "text/plain")
	return w
}

func errorResp(w http.ResponseWriter, msg string, err *error) {
	fmt.Println("Error loading seeds", msg)
	if err != nil {
		fmt.Println(*err)
	}
	http.Error(w, msg, 400)
}

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

	corsHeaders(jsonHeaders(w)).Write([]byte(`{ "message" : "saved" }`))
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

	corsHeaders(jsonHeaders(w)).Write([]byte("{ \"message\" : \"removed\" }"))
	return
}

func Start(host string, port int, s storageCombined, changeCallback func()) {
	srv := RESTServer{
		credStorage: s,
	}

	rtr := mux.NewRouter()

	type x = func(w http.ResponseWriter, req *http.Request)

	wrapper := func(f x) x {
		return func(w http.ResponseWriter, req *http.Request) {
			f(w, req)
			changeCallback()
		}
	}

	// handle CORS OPTIONS
	rtr.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsHeaders(w)
	})

	rtr.HandleFunc("/list", srv.listCredentials).Methods(http.MethodGet)
	rtr.HandleFunc("/add", wrapper(srv.saveCredentials)).Methods(http.MethodPut)
	rtr.HandleFunc("/{id}", wrapper(srv.loadCredentials)).Methods(http.MethodGet)
	rtr.HandleFunc("/{id}", wrapper(srv.removeCredentials)).Methods(http.MethodDelete)

	box := packr.NewBox("../../web/dist/")

	fs := http.FileServer(box)
	rtr.PathPrefix("/").Handler(fs)

	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), rtr)

}
