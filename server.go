package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pingcap/log"
)

type Cake struct {
	ID       string `json: "id"`
	Name     string `json: "name"`
	Picture  string `json: "picture"`
	Category string `json: "category"`
	Comments int    `json: "comments"`
	UpVotes  int    `json: "upvotes"`
}

type PutCake struct {
	Name     string `json: "name"`
	Picture  string `json: "picture"`
	Category string `json: "category"`
	Comments int    `json: "comments"`
	UpVotes  int    `json: "upvotes"`
}

type PatchCake struct {
	Comments int `json: "comments"`
	UpVotes  int `json: "upvotes"`
}

const (
	FAILURE = false
	SUCCESS = true
)

type MetaData struct {
	Count   int    `json: "count"`
	Message string `json: message`
}

type RequestResult struct {
	Success bool     `json: success`
	Meta    MetaData `json: "meta"`
	Data    []Cake   `json: "data"`
}

var cakes []Cake

type apiStaticHandler struct {
	staticPath string
	indexPath  string
}

func (h apiStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	index := filepath.Join(h.staticPath, h.indexPath)
	fi, err := os.Stat(index)
	if !fi.IsDir() {
		http.ServeFile(w, r, index)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createRequestResult(success bool, message string, content []Cake) RequestResult {
	result := RequestResult{Success: success}

	if success {
		result.Meta.Count = len(content)
		result.Meta.Message = message
		result.Data = content
	} else {
		result.Meta.Count = 0
		result.Meta.Message = message
		result.Data = nil
	}

	return result
}

func selectCakeById(cake *Cake, params map[string]string) bool {

	for _, item := range cakes {
		if item.ID == params["id"] {
			*cake = item
			return true
		}
	}

	return false
}

func getAllCakes(w http.ResponseWriter, r *http.Request) {

	var result RequestResult

	if len(cakes) == 0 {
		result = createRequestResult(FAILURE, "The cake catalog is empty", cakes)
	} else {
		result = createRequestResult(SUCCESS, "Showing the cake catalog", cakes)
	}

	data, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getCake(w http.ResponseWriter, r *http.Request) {

	var cake Cake
	var result RequestResult
	var selected []Cake

	if selectCakeById(&cake, mux.Vars(r)) {
		result = createRequestResult(SUCCESS, "Cake was selected", append(selected, cake))
	} else {
		result = createRequestResult(FAILURE, "Cake was not selected", selected)
	}

	data, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func postCake(w http.ResponseWriter, r *http.Request) {

	var cake Cake
	err := json.NewDecoder(r.Body).Decode(&cake)
	if err != nil {
		log.Error(err.Error())
		data, _ := json.Marshal(createRequestResult(FAILURE, err.Error(), nil))
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	cake.ID = uuid.NewString()
	cakes = append(cakes, cake)

	data, _ := json.Marshal(createRequestResult(SUCCESS, "Cake was created", cakes))
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func putCake(w http.ResponseWriter, r *http.Request) {

	var put PutCake
	err := json.NewDecoder(r.Body).Decode(&put)
	if err != nil {
		log.Error(err.Error())
		data, _ := json.Marshal(createRequestResult(FAILURE, err.Error(), nil))
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	var cake Cake
	var result RequestResult
	var puted []Cake

	if selectCakeById(&cake, mux.Vars(r)) {
		for index, item := range cakes {
			if item.ID == cake.ID {
				cake.Name = put.Name
				cake.Picture = put.Picture
				cake.Category = put.Category
				cake.Comments = put.Comments
				cake.UpVotes = put.UpVotes
				cakes[index] = cake
				break
			}
		}
		result = createRequestResult(SUCCESS, "Cake was updated", append(puted, cake))
	} else {
		result = createRequestResult(FAILURE, "Cake was not updated", puted)
	}

	data, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func patchCake(w http.ResponseWriter, r *http.Request) {

	var patch PatchCake
	err := json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		log.Error(err.Error())
		data, _ := json.Marshal(createRequestResult(FAILURE, err.Error(), nil))
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	var cake Cake
	var result RequestResult
	var patched []Cake

	if selectCakeById(&cake, mux.Vars(r)) {
		for index, item := range cakes {
			if item.ID == cake.ID {
				cake.Comments = patch.Comments
				cake.UpVotes = patch.UpVotes
				cakes[index] = cake
				break
			}
		}
		result = createRequestResult(SUCCESS, "Cake was updated", append(patched, cake))
	} else {
		result = createRequestResult(FAILURE, "Cake was not updated", patched)
	}

	data, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func deleteCake(w http.ResponseWriter, r *http.Request) {

	var cake Cake
	var result RequestResult
	var removed []Cake

	if selectCakeById(&cake, mux.Vars(r)) {
		for index, item := range cakes {
			if item.ID == cake.ID {
				cakes = append(cakes[:index], cakes[index+1:]...)
				break
			}
		}
		result = createRequestResult(SUCCESS, "Cake was deleted", append(removed, cake))
	} else {
		result = createRequestResult(FAILURE, "Cake was not deleted", removed)
	}

	data, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {

	usage := apiStaticHandler{staticPath: "static", indexPath: "index.html"}
	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))

	// API Usage
	router.Path("/").Handler(usage)
	router.Path("/api").Handler(usage)
	router.Path("/api/v1").Handler(usage)

	// Cakes Routes
	router.Path("/api/v1/cakes").HandlerFunc(getAllCakes).Methods("GET")
	router.Path("/api/v1/cakes/{id}").HandlerFunc(getCake).Methods("GET")
	router.Path("/api/v1/cakes").HandlerFunc(postCake).Methods("POST")
	router.Path("/api/v1/cakes/{id}").HandlerFunc(putCake).Methods("PUT")
	router.Path("/api/v1/cakes/{id}").HandlerFunc(patchCake).Methods("PATCH")
	router.Path("/api/v1/cakes/{id}").HandlerFunc(deleteCake).Methods("DELETE")

	server := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8086",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
