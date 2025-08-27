package subscription

import (
	"net/http"
	"rest-service/types"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.SubscriptionStore
}

func NewHandler(store types.SubscriptionStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/subscriptions", h.handleGetSubscriptions).Methods(http.MethodGet)
	router.HandleFunc("/subscriptions", h.handleCreateSubscription).Methods(http.MethodPost)
	router.HandleFunc("/subscriptions", h.handleModifySubscription).Methods(http.MethodPut)
	router.HandleFunc("/subscriptions", h.handleDeleteSubscription).Methods(http.MethodDelete)
}

func (h *Handler) handleGetSubscriptions(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) handleModifySubscription(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
}

