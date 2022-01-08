package transport

import "net/http"

func (h *Handler) registerRoutes() {
	h.router.Methods(http.MethodPost).Path("/{typePlural}").HandlerFunc(h.CreateItemHandler())
	h.router.Methods(http.MethodGet).Path("/{typePlural}").HandlerFunc(h.ListItemsHandler())
	h.router.Methods(http.MethodGet).Path("/{typePlural}/{name}").HandlerFunc(h.ReadItemHandler())
	h.router.Methods(http.MethodPut).Path("/{typePlural}/{name}").HandlerFunc(h.ReplaceItemHandler())
	h.router.Methods(http.MethodPatch).Path("/{typePlural}/{name}").HandlerFunc(h.PatchItemHandler())
	h.router.Methods(http.MethodDelete).Path("/{typePlural}/{name}").HandlerFunc(h.DeleteItemHandler())
}
