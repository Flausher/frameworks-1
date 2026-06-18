package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"user.com/m/internal/domain"
	"user.com/m/internal/middleware"
	"user.com/m/internal/service/user"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	Code      string `json:"code"`
	RequestID string `json:"request_id"`
}

func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	requestID := middleware.GetRequestID(r.Context())
	errorResp := ErrorResponse{
		Error:     message,
		Code:      code,
		RequestID: requestID,
	}

	json.NewEncoder(w).Encode(errorResp)
}

type UserHandler struct {
	us *user.UserService
}

func NewUserHandler(us *user.UserService) *UserHandler {
	return &UserHandler{
		us: us,
	}
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format")
		return
	}

	err = uh.us.CreateUser(&user)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		writeError(w, r, http.StatusBadRequest, "MISSING_ID", "User ID is required")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, "INVALID_ID", "Invalid user ID format")
		return
	}

	user, err := uh.us.GetUserById(idInt)
	if err != nil {
		writeError(w, r, http.StatusNotFound, "USER_NOT_FOUND", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := uh.us.GetUsers()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
