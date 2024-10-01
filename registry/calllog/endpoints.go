package calllog

import (
	"emnavisa/webserver/infrastructure/models"
	"emnavisa/webserver/infrastructure/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Handler struct {
	service *Service
}

func newHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (handler *Handler) GetAllCalls(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(utils.UserContextKey).(utils.AuthedUser)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user info")
		return
	}
	var err error
	filterDate := time.Now()
	if dateStr := r.URL.Query().Get("d"); dateStr != "" {
		filterDate, err = time.Parse("2006-01-02", r.URL.Query().Get("d")) // Expected format: YYYY-MM-DD
		if err != nil {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid date format")
			return
		}
	}
	fmt.Println(filterDate)

	// check if can see all or only the his own calls
	//
	calls, err := handler.service.GetAllCallsByAgent(user.Username, filterDate)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve calls")
		return
	}
	// Get Basic Data
	countries := []map[string]string{
		{"name": "Romania"},
		{"name": "Poland"},
		{"name": "Malta"},
		{"name": "Canada"},
		{"name": "France"},
	}
	jobs := []map[string]string{
		{"name": "Constructions"},
		{"name": "Ouvrier général"},
		{"name": "Chef cuisinier"},
		{"name": "Patissier"},
		{"name": "Forgeron"},
		{"name": "Auxiliaire de vie"},
		{"name": "Mecanique auto"},
		{"name": "Couture"},
		{"name": "Technicien de cuir"},
		{"name": "Tapissier"},
		{"name": "Menuisier"},
		{"name": "Soudure plastique"},
		{"name": "Plomberie"},
		{"name": "Peintre de bois"},
		{"name": "Boulangerie"},
		{"name": "Coiffure"},
		{"name": "Agriculteur"},
		{"name": "Babysitter"},
		{"name": "Babysitter"},
		{"name": "Boucher"},
		{"name": "Boucher"},
		{"name": "Ouvriers de nettoyage"},
		{"name": "Livreur"},
		{"name": "AUTRE"},
	}
	// Respond with the call logs as JSON
	utils.RespondWithSuccess(w, map[string]any{
		"calls":     calls,
		"countries": countries,
		"jobs":      jobs,
	})
}

func (handler *Handler) CreateCall(w http.ResponseWriter, r *http.Request) {
	var callLog models.CallLog
	user, ok := r.Context().Value(utils.UserContextKey).(utils.AuthedUser)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user info")
		return
	}

	// Decode the incoming request body into the callLog struct
	err := json.NewDecoder(r.Body).Decode(&callLog)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	callLog.Agent = user.Username
	callLog.CreatedAt = time.Now()
	if err := handler.service.CreateNewCallLog(callLog); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create call log")
		return
	}
	utils.RespondWithSuccess(w, map[string]any{"message": "Call log created successfully"})
}
