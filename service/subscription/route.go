package subscription

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"rest-service/types"
	"rest-service/utils"
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
	log.Println("Started GET subscriptions handler")
	query := r.URL.Query()
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")
	if serviceName != "" && userID != "" {
		log.Printf("Requested subscription with user_id=%s and service_name=%s\n", userID, serviceName)
		sub, err := h.store.GetSubscriptionByUserIDServiceName(userID, serviceName)
		if err != nil {
			log.Println("ERROR: ", err)
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.SubscriptionToSubscriptionResponse(sub))
		log.Println("GET subscription handler finished gracefully")
		return
	}

	if serviceName != "" {
		log.Printf("Requested subscriptions with service_name=%s\n", serviceName)
		subs, err := h.store.GetSubscriptionsByServiceName(serviceName)
		if err != nil {
			log.Println("ERROR: ", err)
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.SubscriptionSliceToSubscriptionResponse(subs))
		log.Println("GET subscription handler finished gracefully")
		return
	}

	if userID != "" {
		log.Printf("Requested subscriptions with user_id=%s\n", userID)
		subs, err := h.store.GetSubscriptionsByUserID(userID)
		if err != nil {
			log.Println("ERROR: ", err)
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		utils.WriteJSON(w, http.StatusOK, utils.SubscriptionSliceToSubscriptionResponse(subs))
		log.Println("GET subscription handler finished gracefully")
		return
	}

	err := fmt.Errorf("no user_id and service_name parameters provided")
	log.Println("ERROR: ", err)
	utils.WriteError(w, http.StatusBadRequest, err)
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
		utils.WriteError(w, http.StatusInternalServerError, err)
		log.Println("ERROR: database operation failed: ", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
	log.Println("CREATE subscription handler finished gracefully")
}

func (h *Handler) handleUpdateSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Started UPDATE subscription handler")

	userID, serviceName, err := utils.ParseValidateUserIDServiceName(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload types.UpdateSubscriptionPayload
	err = utils.ParseAndValidatePayload(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		log.Println("ERROR: ", err)
		return
	}
	log.Printf("Recieved payload is valid, updating database entry with user_id=%s and service_name=%s", userID, serviceName)
	err = h.store.UpdateSubscription(types.Subscription{
		UserID:      userID,
		ServiceName: serviceName,
		Price:       payload.Price,
		StartDate:   time.Time(payload.StartDate),
		EndDate: sql.NullTime{
			Time:  time.Time(payload.EndDate),
			Valid: !time.Time(payload.EndDate).IsZero(),
		},
	})
	if err != nil {
		if err.Error() == "there is no subscription with requested id" {
			utils.WriteError(w, http.StatusNotFound, err)
			log.Println("ERROR: database operation failed: ", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		log.Println("ERROR: database operation failed: ", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
	log.Println("UPDATE subscription handler finished gracefully")
}

func (h *Handler) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Started DELETE subscription handler")
	userID, serviceName, err := utils.ParseValidateUserIDServiceName(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.DeleteSubscription(userID, serviceName)
	if err != nil {
		if err.Error() == "there is no subscription with requested id" {
			utils.WriteError(w, http.StatusNotFound, err)
			log.Println("ERROR: database operation failed: ", err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
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
