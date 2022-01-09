package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nasermirzaei89/core/internal/repository"
)

type Handler struct {
	router   *mux.Router
	itemRepo repository.ItemRepository
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func New(itemRepo repository.ItemRepository) *Handler {
	h := new(Handler)

	h.itemRepo = itemRepo
	h.router = mux.NewRouter()

	h.registerRoutes()

	return h
}
