package account

import (
	"emnavisa/webserver/infrastructure/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler is the http.Handler for this request
type Handler struct {
	service *Service
}

// NewHandler will create a new Handler to handle this request
func newHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// Handle will handle the incoming request
type UserLoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (handler *Handler) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(utils.UserContextKey).(utils.AuthedUser)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user info")
		return
	}

	// Return the user information as a JSON response
	utils.RespondWithSuccess(w, map[string]any{
		"username": user.Username,
		"r":        user.Role,
	})
}

func (handler *Handler) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	var creds UserLoginCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil || creds.Username == "" || creds.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = handler.service.UserCreate(creds.Username, creds.Password)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}
	utils.RespondWithSuccess(w, "created successfully")
}

func (handler *Handler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var creds UserLoginCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil || creds.Username == "" || creds.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	userID, err := handler.service.Authenticate(creds.Username, creds.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}

	// Store the access token in the database
	accessToken, err := handler.service.StoreAccessToken(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to store access token")
		return
	}

	utils.RespondWithSuccess(w, accessToken)
}
