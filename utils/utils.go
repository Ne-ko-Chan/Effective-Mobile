package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-service/types"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, types.Error{
		Error: err.Error(),
	})
}

func ParseAndValidatePayload(r *http.Request, payload any) error {
	if err := ParseJSON(r, &payload); err != nil {
		return fmt.Errorf("ERROR: payload parsing failed: %v", err)
	}

	if err := Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return fmt.Errorf("ERROR: payload is invalid: %v", errors)
	}

	return nil
}

func SubscriptionToSubscriptionResponse(s *types.Subscription) types.SubscriptionResponse {
	if s.EndDate.Valid {
		return types.SubscriptionResponseWithEnd{
			UserID:      s.UserID,
			ServiceName: s.ServiceName,
			Price:       s.Price,
			StartDate:   types.CustomTime(s.StartDate),
			EndDate:     types.CustomTime(s.EndDate.Time),
		}
	} else {
		return types.SubscriptionResponseNoEnd{
			UserID:      s.UserID,
			ServiceName: s.ServiceName,
			Price:       s.Price,
			StartDate:   types.CustomTime(s.StartDate),
		}
	}
}

func SubscriptionSliceToSubscriptionResponse(subs []types.Subscription) []types.SubscriptionResponse {
	res := make([]types.SubscriptionResponse, len(subs))
	for i := range res {
		res[i] = SubscriptionToSubscriptionResponse(&subs[i])
	}
	return res
}


func ParseValidateUserIDServiceName(r *http.Request) (string, string, error) {
	query := r.URL.Query()
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")
	if userID == "" {
		return "", "", fmt.Errorf("no user_id query parameter provided")
	}
	if serviceName == "" {
		return "", "", fmt.Errorf("no service_name query parameter provided")
	}

	return userID, serviceName, nil
}
