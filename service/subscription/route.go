package subscription

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"rest-service/types"
	"rest-service/utils"
	"strconv"
	"time"

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
	router.HandleFunc("/subscriptions", h.handleUpdateSubscription).Methods(http.MethodPut)
	router.HandleFunc("/subscriptions", h.handleDeleteSubscription).Methods(http.MethodDelete)
	router.HandleFunc("/subscriptions/sum", h.handleGetSum).Methods(http.MethodGet)
}

func (h *Handler) handleGetSubscriptions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	listFlag := query.Get("list")
	if listFlag == "" || listFlag == "false" {
		h.handleReadSubscription(w, r)
		return
	}
	if listFlag == "true" {
		h.handleListSubscriptions(w, r)
		return
	}
	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect value of list parameter. Must be bool if present. If not present, defaults to 'false'"))
}

func (h *Handler) handleReadSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Started READ subscription handler")
	query := r.URL.Query()
	subscriptionID := query.Get("id")
	if subscriptionID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no id query parameter specified"))
		log.Println("ERROR: no id query parameter specified")
		return
	}
	id, err := strconv.ParseInt(subscriptionID, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: ", err)
		return
	}
	log.Println("Requested subscription id=",id)
	sub, err := h.store.GetSubscriptionByID(int(id))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: ", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.SubscriptionToSubscriptionResponse(sub))
	log.Println("READ subscription handler finished gracefully")
}

func (h *Handler) handleListSubscriptions(w http.ResponseWriter, _ *http.Request) {
	log.Println("Started LIST subscriptions handler")
	subs, err := h.store.GetSubscriptions()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: ", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.SubscriptionSliceToSubscriptionResponse(subs))
	log.Println("LIST subscriptions handler finished gracefully")
}

func (h *Handler) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Started CREATE subscription handler")

	var payload types.CreateSubscriptionPayload
	err := utils.ParseAndValidatePayload(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: ", err)
		return
	}
	log.Println("Recieved payload is valid, creating database entry")

	err = h.store.CreateSubscription(types.Subscription{
		UserID:      payload.UserID,
		ServiceName: payload.ServiceName,
		Price:       payload.Price,
		StartDate:   time.Time(payload.StartDate),
		EndDate: sql.NullTime{
			Time:  time.Time(payload.EndDate),
			Valid: !time.Time(payload.EndDate).IsZero(),
		},
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: database operation failed: ", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
	log.Println("CREATE subscription handler finished gracefully")
}

func (h *Handler) handleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Started UPDATE subscription handler")

	var payload types.UpdateSubscriptionPayload
	err := utils.ParseAndValidatePayload(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: ", err)
		return
	}
	log.Println("Recieved payload is valid, updating database entry")
	err = h.store.UpdateSubscription(types.Subscription{
		ID:          payload.ID,
		UserID:      payload.UserID,
		ServiceName: payload.ServiceName,
		Price:       payload.Price,
		StartDate:   time.Time(payload.StartDate),
		EndDate: sql.NullTime{
			Time:  time.Time(payload.EndDate),
			Valid: !time.Time(payload.EndDate).IsZero(),
		},
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: database operation failed: ", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
	log.Println("UPDATE subscription handler finished gracefully")
}

func (h *Handler) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Started DELETE subscription handler")

	var payload types.DeleteSubscriptionPayload
	err := utils.ParseAndValidatePayload(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: ", err)
		return
	}
	log.Println("Recieved payload is valid, updating database entry")

	err = h.store.DeleteSubscription(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: database operation failed: ", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
	log.Println("DELETE subscription handler finished gracefully")
}

func (h *Handler) handleGetSum(w http.ResponseWriter, r *http.Request) {
	log.Println("Started GET SUM handler")
	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")

	sum, err := h.store.GetSumPeriod(from, to, userID, serviceName)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: ", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]int{"sum": sum})
	log.Println("GET SUM handler finished gracefully")
}
