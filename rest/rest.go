package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/user"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/jdevelop/passkeeper"
	"github.com/jdevelop/passkeeper/storage"
)

type storageCombined interface {
	storage.SeedStorageRead
	storage.SeedStorageWrite
	storage.SeedStorageList
	storage.SeedStorageRemove
}

type RESTServer struct {
	SeedStorage storageCombined
}

func corsHeaders(w http.ResponseWriter) http.ResponseWriter {
	hdr := w.Header()
	hdr.Set("Access-Control-Allow-Origin", "*")
	hdr.Set("Access-Control-Allow-Methods", "GET,PUT,DELETE,POST")
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

func (r *RESTServer) listSeeds(w http.ResponseWriter, _ *http.Request) {
	seeds, err := r.SeedStorage.ListSeeds()
	if err != nil {
		errorResp(w, "Can't load seeds", &err)
		return
	}
	data, err := json.Marshal(seeds)
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

func (r *RESTServer) loadSeed(w http.ResponseWriter, req *http.Request) {
	seedId, ok := req.URL.Query()["seed"]
	if !ok || len(seedId) != 1 {
		errorResp(w, "No required parameter", nil)
		return
	}
	seed, err := r.SeedStorage.LoadSeed(seedId[0])

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

	return
}

func (r *RESTServer) saveSeed(w http.ResponseWriter, req *http.Request) {

	data, err := ioutil.ReadAll(req.Body)

	if err != nil {
		errorResp(w, "Can't read seed object", &err)
		return
	}

	var s passkeeper.Seed

	if err = json.Unmarshal(data, &s); err != nil {
		errorResp(w, "Can't unmarshal seed object", &err)
		return
	}

	if err = r.SeedStorage.SaveSeed(s); err != nil {
		errorResp(w, "Can't save seed", &err)
		return
	}

	corsHeaders(jsonHeaders(w)).Write([]byte("{ \"message\" : \"saved\" }"))
	return
}

func (r *RESTServer) removeSeed(w http.ResponseWriter, req *http.Request) {

	seedId, ok := req.URL.Query()["seed"]
	if !ok || len(seedId) != 1 {
		errorResp(w, "No required parameter", nil)
		return
	}

	if err := r.SeedStorage.RemoveSeed(seedId[0]); err != nil {
		errorResp(w, "Can't remove seed", &err)
		return
	}

	corsHeaders(jsonHeaders(w)).Write([]byte("{ \"message\" : \"removed\" }"))
	return
}

func Start(host string, port int, s storageCombined, changeCallback func()) {
	srv := RESTServer{
		SeedStorage: s,
	}

	rtr := mux.NewRouter()

	type x = func(w http.ResponseWriter, req *http.Request)

	wrapper := func(f x) x {
		return func(w http.ResponseWriter, req *http.Request) {
			f(w, req)
			changeCallback()
		}
	}

	rtr.HandleFunc("/api/seeds", srv.listSeeds).Methods("GET")
	rtr.HandleFunc("/api/seed", wrapper(srv.saveSeed)).Methods("PUT")
	rtr.HandleFunc("/api/seed", wrapper(srv.loadSeed)).Methods("GET")
	rtr.HandleFunc("/api/seed", wrapper(srv.removeSeed)).Methods("DELETE")

	var staticPath string

	u, err := user.Current()
	if err != nil {
		staticPath = "/root/web"
	} else {
		staticPath = filepath.Join(u.HomeDir, "web")
	}

	fs := http.FileServer(http.Dir(staticPath))
	rtr.PathPrefix("/").Handler(fs)

	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), rtr)

}
