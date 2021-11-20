package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nasermirzaei89/env"
	"github.com/pkg/errors"
)

func main() {
	repo := NewMemoryRepository()
	router := mux.NewRouter()

	router.Methods(http.MethodPost).Path("/{typePlural}").HandlerFunc(CreateItemHandler(repo))
	router.Methods(http.MethodGet).Path("/{typePlural}").HandlerFunc(ListItemsHandler(repo))
	router.Methods(http.MethodGet).Path("/{typePlural}/{name}").HandlerFunc(ReadItemHandler(repo))
	router.Methods(http.MethodPut).Path("/{typePlural}/{name}").HandlerFunc(ReplaceItemHandler(repo))
	router.Methods(http.MethodDelete).Path("/{typePlural}/{name}").HandlerFunc(DeleteItemHandler(repo))

	err := http.ListenAndServe(env.GetString("API_ADDRESS", ":80"), router)
	if err != nil {
		panic(errors.Wrap(err, "error on listen and serve http"))
	}
}
