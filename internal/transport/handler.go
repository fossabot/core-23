package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nasermirzaei89/core/internal/repository"
)

type Handler struct {
	router   *mux.Router
	itemRepo repository.Item
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func New(itemRepo repository.Item) *Handler {
	h := new(Handler)

	h.itemRepo = itemRepo
	h.router = mux.NewRouter()

	h.registerRoutes()

	return h
}
