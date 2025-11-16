package handlers

import (
	"encoding/json"
	"net/http"
	"reviewer-assignment-service/internal/app/response_errors"
	"reviewer-assignment-service/internal/app/transport/dtos"
	"reviewer-assignment-service/internal/app/transport/mappers"
	"reviewer-assignment-service/internal/app/validators"
	"strconv"

	"reviewer-assignment-service/internal/domain/services"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService services.UserService
	prService   services.PullRequestService
}

func NewUserHandler(userService services.UserService, prService services.PullRequestService) *UserHandler {
	return &UserHandler{
		userService: userService,
		prService:   prService,
	}
}

func (h *UserHandler) SetUserActive(w http.ResponseWriter, r *http.Request) {
	var req dtos.SetUserActiveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response_errors.SendError(w, "INVALID_REQUEST", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateSetUserActiveRequest(&req); err != nil {
		response_errors.SendError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}
	userID, _ := strconv.Atoi(req.UserID)

	_, err := h.userService.GetByID(userID)
	if err != nil {
		response_errors.SendError(w, "NOT_FOUND", "User not found", http.StatusNotFound)
		return
	}

	if err := h.userService.SetActive(userID, req.IsActive); err != nil {
		response_errors.SendError(w, "INTERNAL_ERROR", "Failed to update user", http.StatusInternalServerError)
		return
	}

	updatedUser, err := h.userService.GetByID(userID)
	if err != nil {
		response_errors.SendError(w, "INTERNAL_ERROR", "Failed to get updated user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user": mappers.UserToResponse(updatedUser),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUserReviewPRs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	if err := validators.ValidateUserID(userID); err != nil {
		response_errors.SendError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	userIDInt, _ := strconv.Atoi(userID)

	_, err := h.userService.GetByID(userIDInt)
	if err != nil {
		response_errors.SendError(w, "NOT_FOUND", "User not found", http.StatusNotFound)
		return
	}

	prs, err := h.prService.GetByReviewerID(userIDInt)
	if err != nil {
		response_errors.SendError(w, "INTERNAL_ERROR", "Failed to get user PRs", http.StatusInternalServerError)
		return
	}

	response := mappers.UserToResponseWithPRs(userID, prs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response_errors.SendError(w, "INVALID_REQUEST", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateCreateUserRequest(&req); err != nil {
		response_errors.SendError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	user := mappers.CreateUserRequestToDomain(req)

	if err := h.userService.Create(user); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	createdUser, err := h.userService.GetByID(user.ID)
	if err != nil {
		response_errors.SendError(w, "INTERNAL_ERROR", "Failed to get created user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user": mappers.UserToDetailedResponse(createdUser),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	var req dtos.DeactivateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response_errors.SendError(w, "INVALID_REQUEST", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateDeactivateUserRequest(&req); err != nil {
		response_errors.SendError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(req.UserID)

	_, err := h.userService.GetByID(userID)
	if err != nil {
		response_errors.SendError(w, "NOT_FOUND", "User not found", http.StatusNotFound)
		return
	}

	if err := h.userService.Deactivate(userID); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	updatedUser, err := h.userService.GetByID(userID)
	if err != nil {
		response_errors.SendError(w, "INTERNAL_ERROR", "Failed to get updated user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user": mappers.UserToDetailedResponse(updatedUser),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll()
	if err != nil {
		response_errors.SendError(w, "INTERNAL_ERROR", "Failed to get users", http.StatusInternalServerError)
		return
	}

	var userResponses []dtos.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, mappers.UserToDetailedResponse(user))
	}

	response := map[string]interface{}{
		"users": userResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")

	if err := validators.ValidateUserID(userIDStr); err != nil {
		response_errors.SendError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)
	user, err := h.userService.GetByID(userID)
	if err != nil {
		response_errors.SendError(w, "NOT_FOUND", "User not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"user": mappers.UserToResponse(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if err := validators.ValidateEmail(email); err != nil {
		response_errors.SendError(w, "INVALID_REQUEST", err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByEmail(email)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := map[string]interface{}{
		"user": mappers.UserToDetailedResponse(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
