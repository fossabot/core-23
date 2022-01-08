package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nasermirzaei89/core/internal/core"
	"github.com/nasermirzaei89/core/internal/repository"
	"github.com/pkg/errors"
)

type HTTPError struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (h *Handler) CreateItemHandler() http.HandlerFunc {
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
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type field is not valid, it should an string that matches the regex '%s'", core.TypeRegex)})

			return
		}

		var req CreateItemRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on decode request body", Error: err.Error()})

			return
		}

		if !isValidName(req.Name) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("name field is not valid, it should an string that matches the regex '%s'", core.NameRegex)})

			return
		}

		_, err = h.itemRepo.FindByTypeAndName(r.Context(), typ, req.Name)
		if err != nil {
			if !errors.Is(err, repository.ErrItemNotFound) {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by type and name from the repository", Error: err.Error()})

				return
			}
		} else {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("%s with name '%s' already exists", typ, req.Name)})

			return
		}

		now := time.Now()

		item := core.Item{
			UUID:      uuid.NewString(),
			Type:      typ,
			Name:      req.Name,
			Data:      req.Data,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = h.itemRepo.Insert(r.Context(), item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on insert item to the repository", Error: err.Error()})

			return
		}

		rsp := ItemFromEntity(item)

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(rsp)
	}
}

func isValidName(name string) bool {
	return regexp.MustCompile(core.NameRegex).MatchString(name)
}

func isValidType(typ string) bool {
	return regexp.MustCompile(core.TypeRegex).MatchString(typ)
}

func (h *Handler) ListItemsHandler() http.HandlerFunc {
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
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type field is not valid, it should an string that matches the regex '%s'", core.TypeRegex)})

			return
		}

		items, err := h.itemRepo.ListByType(r.Context(), typ)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on list items by type from the repository", Error: err.Error()})

			return
		}

		rsp := ListItemsResponse{Items: make([]Item, len(items))}

		for i := range items {
			rsp.Items[i] = ItemFromEntity(items[i])
		}

		_ = json.NewEncoder(w).Encode(rsp)
	}
}

func (h *Handler) ReadItemHandler() http.HandlerFunc {
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
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type field is not valid, it should an string that matches the regex '%s'", core.TypeRegex)})

			return
		}

		name := mux.Vars(r)["name"]

		item, err := h.itemRepo.FindByTypeAndName(r.Context(), typ, name)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrItemNotFound):
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("%s with name '%s' not found", typ, name)})
			default:
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by type and name from the repository", Error: err.Error()})
			}

			return
		}

		rsp := ItemFromEntity(item)

		_ = json.NewEncoder(w).Encode(rsp)
	}
}

func (h *Handler) ReplaceItemHandler() http.HandlerFunc {
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
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type parameter is not valid, it should an string that matches the regex '%s'", core.TypeRegex)})

			return
		}

		name := mux.Vars(r)["name"]

		if !isValidName(name) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("name parameter is not valid, it should an string that matches the regex '%s'", core.NameRegex)})

			return
		}

		var req ReplaceItemRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on decode request body", Error: err.Error()})

			return
		}

		item, err := h.itemRepo.FindByTypeAndName(r.Context(), typ, name)
		if err != nil {
			if errors.Is(err, repository.ErrItemNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("%s with name '%s' not found", typ, name)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by type and name from the repository", Error: err.Error()})
			}

			return
		}

		item.Data = req.Data
		item.UpdatedAt = time.Now()

		err = h.itemRepo.Replace(r.Context(), item.UUID, item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on replace item in the repository", Error: err.Error()})

			return
		}

		rsp := ItemFromEntity(item)

		_ = json.NewEncoder(w).Encode(rsp)
	}
}

func (h *Handler) DeleteItemHandler() http.HandlerFunc {
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
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("type parameter is not valid, it should an string that matches the regex '%s'", core.TypeRegex)})

			return
		}

		name := mux.Vars(r)["name"]

		if !isValidName(name) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("name parameter is not valid, it should an string that matches the regex '%s'", core.NameRegex)})

			return
		}

		item, err := h.itemRepo.FindByTypeAndName(r.Context(), typ, name)
		if err != nil {
			if errors.Is(err, repository.ErrItemNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: fmt.Sprintf("%s with name '%s' not found", typ, name)})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on find item by type and name from the repository", Error: err.Error()})
			}

			return
		}

		err = h.itemRepo.Delete(r.Context(), item.UUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(HTTPError{Message: "error on delete item from the repository", Error: err.Error()})

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
