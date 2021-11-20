package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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

		iName, ok := item[ItemFieldName]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "name field is not provided"})

			return
		}

		name, ok := iName.(string)
		if !ok || !isValidName(name) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("name field is not valid, it should an string that matches the regex '%s'", NameRegex)})

			return
		}

		_, err = repo.FindByName(r.Context(), name)
		if err != nil {
			if !errors.Is(err, ErrItemNotFound) {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by name from the repository", Error: err.Error()})

				return
			}
		} else {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("item with name '%s' already exists", name)})

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

func isValidName(name string) bool {
	return regexp.MustCompile(NameRegex).MatchString(name)
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
		name := mux.Vars(r)["name"]

		item, err := repo.FindByName(r.Context(), name)
		if err != nil {
			if errors.Is(err, ErrItemNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("item with name '%s' not found", name)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by name from the repository", Error: err.Error()})
			}

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

		name := mux.Vars(r)["name"]

		item, err := repo.FindByName(r.Context(), name)
		if err != nil {
			if errors.Is(err, ErrItemNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("item with name '%s' not found", name)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by name from the repository", Error: err.Error()})
			}

			return
		}

		createdAt := item[ItemFieldCreatedAt]
		itemUUID := item[ItemFieldUUID].(string)

		item = newItem
		item[ItemFieldName] = name
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
		name := mux.Vars(r)["name"]

		item, err := repo.FindByName(r.Context(), name)
		if err != nil {
			if errors.Is(err, ErrItemNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("item with name '%s' not found", name)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by name from the repository", Error: err.Error()})
			}

			return
		}

		err = repo.Delete(r.Context(), item[ItemFieldUUID].(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on delete item from the repository", Error: err.Error()})

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
