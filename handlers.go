package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gertd/go-pluralize"
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
		pc := pluralize.NewClient()

		typePlural := mux.Vars(r)["typePlural"]

		if !pc.IsPlural(typePlural) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "you should set plural form of the type"})

			return
		}

		typ := pc.Singular(typePlural)

		if !isValidType(typ) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type field is not valid, it should an string that matches the regex '%s'", TypeRegex)})

			return
		}

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

		_, err = repo.FindByTypeAndName(r.Context(), typ, name)
		if err != nil {
			if !errors.Is(err, ErrItemNotFound) {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by type and name from the repository", Error: err.Error()})

				return
			}
		} else {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("%s with name '%s' already exists", typ, name)})

			return
		}

		item[ItemFieldUUID] = uuid.NewString()
		item[ItemFieldType] = typ

		now := time.Now().Format(time.RFC3339)
		item[ItemFieldCreatedAt] = now
		item[ItemFieldUpdatedAt] = now

		err = repo.Insert(r.Context(), Item(item))
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

func isValidType(typ string) bool {
	return regexp.MustCompile(TypeRegex).MatchString(typ)
}

func ListItemsHandler(repo Repository) http.HandlerFunc {
	type Response struct {
		Items []Item `json:"items"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		pc := pluralize.NewClient()

		typePlural := mux.Vars(r)["typePlural"]

		if !pc.IsPlural(typePlural) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "you should set plural form of the type"})

			return
		}

		typ := pc.Singular(typePlural)

		if !isValidType(typ) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type field is not valid, it should an string that matches the regex '%s'", TypeRegex)})

			return
		}

		items, err := repo.ListByType(r.Context(), typ)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on list items by type from the repository", Error: err.Error()})

			return
		}

		rsp := Response{Items: items}

		_ = json.NewEncoder(w).Encode(rsp)
	}
}

func ReadItemHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pc := pluralize.NewClient()

		typePlural := mux.Vars(r)["typePlural"]

		if !pc.IsPlural(typePlural) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "you should set plural form of the type"})

			return
		}

		typ := pc.Singular(typePlural)

		if !isValidType(typ) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type field is not valid, it should an string that matches the regex '%s'", TypeRegex)})

			return
		}

		name := mux.Vars(r)["name"]

		item, err := repo.FindByTypeAndName(r.Context(), typ, name)
		if err != nil {
			if errors.Is(err, ErrItemNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("%s with name '%s' not found", typ, name)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by type and name from the repository", Error: err.Error()})
			}

			return
		}

		_ = json.NewEncoder(w).Encode(item)
	}
}

func ReplaceItemHandler(repo Repository) http.HandlerFunc {
	type Request map[string]interface{}

	return func(w http.ResponseWriter, r *http.Request) {
		pc := pluralize.NewClient()

		typePlural := mux.Vars(r)["typePlural"]

		if !pc.IsPlural(typePlural) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "you should set plural form of the type"})

			return
		}

		typ := pc.Singular(typePlural)

		if !isValidType(typ) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type field is not valid, it should an string that matches the regex '%s'", TypeRegex)})

			return
		}

		var newItem Request

		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on decode request body", Error: err.Error()})

			return
		}

		name := mux.Vars(r)["name"]

		item, err := repo.FindByTypeAndName(r.Context(), typ, name)
		if err != nil {
			if errors.Is(err, ErrItemNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("%s with name '%s' not found", typ, name)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by type and name from the repository", Error: err.Error()})
			}

			return
		}

		createdAt := item[ItemFieldCreatedAt]
		itemUUID := item.UUID()

		item = Item(newItem)
		item[ItemFieldUUID] = itemUUID
		item[ItemFieldName] = name
		item[ItemFieldType] = typ
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
		pc := pluralize.NewClient()

		typePlural := mux.Vars(r)["typePlural"]

		if !pc.IsPlural(typePlural) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "you should set plural form of the type"})

			return
		}

		typ := pc.Singular(typePlural)

		if !isValidType(typ) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type field is not valid, it should an string that matches the regex '%s'", TypeRegex)})

			return
		}

		name := mux.Vars(r)["name"]

		item, err := repo.FindByTypeAndName(r.Context(), typ, name)
		if err != nil {
			if errors.Is(err, ErrItemNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("%s with name '%s' not found", typ, name)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by type and name from the repository", Error: err.Error()})
			}

			return
		}

		err = repo.Delete(r.Context(), item.UUID())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on delete item from the repository", Error: err.Error()})

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
