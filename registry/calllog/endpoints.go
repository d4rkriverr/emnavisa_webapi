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
	if user.Role != "ADMN" {
		if user.Role != "CACE" {
			utils.RespondWithError(w, http.StatusUnauthorized, "unauthorized access")
			return
		}
	}

	var err error
	filterDate := time.Now()
	if dateStr := r.URL.Query().Get("d"); dateStr != "" {
		filterDate, err = time.Parse("2006-01-02", r.URL.Query().Get("d")) // Expected format: YYYY-MM-DD
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid date format")
			return
		}
	}

	// check if can see all or only the his own calls
	var calls []models.CallLog
	if user.Role == "ADMN" {
		calls, err = handler.service.GetAllCalls(filterDate)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve calls")
			return
		}
	} else {
		calls, err = handler.service.GetAllCallsByAgent(user.Username, filterDate)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve calls")
			return
		}
	}
	//

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
		{"name": "Boucherie"},
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

func (handler *Handler) EditCall(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(utils.UserContextKey).(utils.AuthedUser)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user info")
		return
	}
	var updatedCallLog models.CallLog
	if err := json.NewDecoder(r.Body).Decode(&updatedCallLog); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := handler.service.UpdateCallLog(user.Username, updatedCallLog)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to update call log", http.StatusInternalServerError)
		return
	}
	utils.RespondWithSuccess(w, map[string]any{"message": "Call log updated successfully"})
}

func (handler *Handler) DeleteCall(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(utils.UserContextKey).(utils.AuthedUser)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user info")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := handler.service.DeleteCallLog(user.Username, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to delete call log", http.StatusInternalServerError)
		return
	}
	utils.RespondWithSuccess(w, map[string]any{"message": "Call log deleted successfully"})
}
