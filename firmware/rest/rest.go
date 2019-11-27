package rest

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/jdevelop/passkeeper/firmware/pass"
	"github.com/jdevelop/passkeeper/firmware/storage"
)

type storageCombined interface {
	storage.CredentialsStorageList
	storage.CredentialsStorageRead
	storage.CredentialsStorageRemove
	storage.CredentialsStorageWrite
	storage.CredentialsStorageBackup
	storage.CredentialsStorageRestore
	storage.MutableKey
}

type RESTServer struct {
	credStorage        storageCombined
	passwordGen        pass.PasswordGenerator
	cardPasswordChange func([]byte) error
}

func corsHeaders(w http.ResponseWriter) http.ResponseWriter {
	hdr := w.Header()
	hdr.Set("Access-Control-Allow-Origin", "*")
	hdr.Set("Access-Control-Allow-Methods", "OPTIONS,GET,PUT,DELETE,POST")
	hdr.Set("Access-Control-Allow-Headers", "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
	return w
}

var (
	removed  = []byte(`{ "message" : "removed" }`)
	saved    = []byte(`{ "message" : "saved" }`)
	restored = []byte(`{ "message" : "restored" }`)
)

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

func Start(host string, port int, s storageCombined, pwdGen pass.PasswordGenerator,
	cardAccess func([]byte) error,
	changeCallback func()) {
	srv := RESTServer{
		credStorage:        s,
		passwordGen:        pwdGen,
		cardPasswordChange: cardAccess,
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

	rtr.HandleFunc("/backup", srv.backupCredentials).Methods(http.MethodGet)
	rtr.HandleFunc("/restore", wrapper(srv.restoreCredentials)).Methods(http.MethodPost)
	rtr.HandleFunc("/cardpassword", wrapper(srv.setPassword)).Methods(http.MethodPut)
	rtr.HandleFunc("/list", srv.listCredentials).Methods(http.MethodGet)
	rtr.HandleFunc("/add", wrapper(srv.saveCredentials)).Methods(http.MethodPut)
	rtr.HandleFunc("/generate", wrapper(srv.generatePasswords)).Methods(http.MethodGet)
	rtr.HandleFunc("/{id}", wrapper(srv.loadCredentials)).Methods(http.MethodGet)
	rtr.HandleFunc("/{id}", wrapper(srv.removeCredentials)).Methods(http.MethodDelete)

	box := packr.NewBox("../../web/")

	fs := http.FileServer(box)
	rtr.PathPrefix("/").Handler(fs)

	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), rtr)

}
