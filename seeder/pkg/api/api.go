package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/odio4u/memstore/seeder/pkg/memstore"
)

type Api struct {
	memstore *memstore.MemStore
}

func NewApi(memstore *memstore.MemStore) *Api {
	return &Api{
		memstore: memstore,
	}
}

func SetRoutes(router *mux.Router, api *Api) {
	router.HandleFunc("/seeder", api.SeederView).Methods("GET")
}

func (a *Api) SeederView(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Query())
	region := r.URL.Query().Get("region")

	seederData := a.memstore.GetSeeders(region)
	response, _ := json.Marshal(seederData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}
