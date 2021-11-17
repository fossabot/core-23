package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPError struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func CreateItemHandler(repo Repository) http.HandlerFunc {
	type Request map[string]interface{}

	return func(w http.ResponseWriter, r *http.Request) {
		var item Request

		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on decode request body", Error: err.Error()})

			return
		}

		item[ItemFieldUUID] = uuid.NewString()

		now := time.Now().Format(time.RFC3339)
		item[ItemFieldCreatedAt] = now
		item[ItemFieldUpdatedAt] = now

		err = repo.Insert(r.Context(), item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on insert item to the repository", Error: err.Error()})

			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(item)
	}
}

func ListItemsHandler(repo Repository) http.HandlerFunc {
	type Response struct {
		Items []map[string]interface{} `json:"items"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		items, err := repo.List(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on list items from the repository", Error: err.Error()})

			return
		}

		rsp := Response{Items: items}

		_ = json.NewEncoder(w).Encode(rsp)
	}
}

func ReadItemHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemUUID := mux.Vars(r)["itemUUID"]

		item, err := repo.Find(r.Context(), itemUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item from the repository", Error: err.Error()})

			return
		}

		if item == nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item from the repository", Error: err.Error()})

			return
		}

		_ = json.NewEncoder(w).Encode(item)
	}
}

func ReplaceItemHandler(repo Repository) http.HandlerFunc {
	type Request map[string]interface{}

	return func(w http.ResponseWriter, r *http.Request) {
		var newItem Request

		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on decode request body", Error: err.Error()})

			return
		}

		itemUUID := mux.Vars(r)["itemUUID"]

		item, err := repo.Find(r.Context(), itemUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item from the repository", Error: err.Error()})

			return
		}

		if item == nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item from the repository", Error: err.Error()})

			return
		}

		createdAt := item[ItemFieldCreatedAt]

		item = newItem
		item[ItemFieldUUID] = itemUUID
		item[ItemFieldCreatedAt] = createdAt
		item[ItemFieldUpdatedAt] = time.Now().Format(time.RFC3339)

		err = repo.Replace(r.Context(), itemUUID, item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on replace item in the repository", Error: err.Error()})

			return
		}

		_ = json.NewEncoder(w).Encode(item)
	}
}

func DeleteItemHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemUUID := mux.Vars(r)["itemUUID"]

		item, err := repo.Find(r.Context(), itemUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item from the repository", Error: err.Error()})

			return
		}

		if item == nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item from the repository", Error: err.Error()})

			return
		}

		err = repo.Delete(r.Context(), itemUUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on delete item from the repository", Error: err.Error()})

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
