package handlers

import (
	"encoding/json"
	"net/http"
	"reviewer-assignment-service/internal/app/response_errors"
	"reviewer-assignment-service/internal/app/transport/dtos"
	"reviewer-assignment-service/internal/app/transport/mappers"
	"reviewer-assignment-service/internal/app/validators"
	"reviewer-assignment-service/internal/domain/models"
	"reviewer-assignment-service/internal/domain/services"
	"time"

	"github.com/go-chi/chi/v5"
)

type PullRequestHandler struct {
	prService   services.PullRequestService
	userService services.UserService
}

func NewPullRequestHandler(prService services.PullRequestService, userService services.UserService) *PullRequestHandler {
	return &PullRequestHandler{
		prService:   prService,
		userService: userService,
	}
}

func (h *PullRequestHandler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreatePullRequestRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response_errors.SendError(w, "INVALID_JSON", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateCreatePullRequestRequest(&req); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	author, err := h.userService.GetByID(req.AuthorID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	pr := &models.PullRequest{
		Name:      req.Name,
		Status:    models.StatusOpen,
		Author:    author,
		Reviewers: make([]*models.User, 0),
		CreatedAt: time.Now(),
	}

	for _, reviewerID := range req.Reviewers {
		reviewer, err := h.userService.GetByID(reviewerID)
		if err != nil {
			response_errors.HandleServiceError(w, err)
			return
		}
		if err := pr.AddReviewer(reviewer); err != nil {
			response_errors.HandleServiceError(w, err)
			return
		}
	}

	if err := h.prService.Create(pr); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToPullRequestResponse(pr)
	sendJSONResponse(w, http.StatusCreated, response)
}

func (h *PullRequestHandler) GetPullRequestByID(w http.ResponseWriter, r *http.Request) {
	prIDStr := chi.URLParam(r, "id")
	prID, err := validators.ValidatePullRequestID(prIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	pr, err := h.prService.GetByID(prID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToPullRequestResponse(pr)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *PullRequestHandler) GetPullRequestsByAuthor(w http.ResponseWriter, r *http.Request) {
	authorIDStr := chi.URLParam(r, "authorID")
	authorID, err := validators.ValidateAuthorID(authorIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	prs, err := h.prService.GetByAuthorID(authorID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToPullRequestListResponse(prs)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *PullRequestHandler) GetPullRequestsByReviewer(w http.ResponseWriter, r *http.Request) {
	reviewerIDStr := chi.URLParam(r, "reviewerID")
	reviewerID, err := validators.ValidateReviewerID(reviewerIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	prs, err := h.prService.GetByReviewerID(reviewerID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToPullRequestListResponse(prs)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *PullRequestHandler) UpdatePullRequest(w http.ResponseWriter, r *http.Request) {
	prIDStr := chi.URLParam(r, "id")
	prID, err := validators.ValidatePullRequestID(prIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	var req dtos.UpdatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response_errors.SendError(w, "INVALID_JSON", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateUpdatePullRequestRequest(&req); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	pr, err := h.prService.GetByID(prID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	pr.Name = req.Name
	pr.Status = models.PRStatus(req.Status)

	pr.Reviewers = make([]*models.User, 0)
	for _, reviewerID := range req.Reviewers {
		reviewer, err := h.userService.GetByID(reviewerID)
		if err != nil {
			response_errors.HandleServiceError(w, err)
			return
		}
		if err := pr.AddReviewer(reviewer); err != nil {
			response_errors.HandleServiceError(w, err)
			return
		}
	}

	if err := h.prService.Update(pr); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToPullRequestResponse(pr)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *PullRequestHandler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	prIDStr := chi.URLParam(r, "id")
	prID, err := validators.ValidatePullRequestID(prIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	pr, err := h.prService.GetByID(prID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	if err := h.prService.MergeRequest(pr); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToPullRequestResponse(pr)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *PullRequestHandler) ReassignReviewers(w http.ResponseWriter, r *http.Request) {
	prIDStr := chi.URLParam(r, "id")
	prID, err := validators.ValidatePullRequestID(prIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	var req dtos.ReassignReviewersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response_errors.SendError(w, "INVALID_JSON", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateReassignReviewersRequest(&req); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	pr, err := h.prService.GetByID(prID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	oldReviewer, err := h.userService.GetByID(req.OldReviewerID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	if err := h.prService.ReassignReviewers(pr, oldReviewer); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	updatedPR, err := h.prService.GetByID(prID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToPullRequestResponse(updatedPR)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *PullRequestHandler) AddReviewer(w http.ResponseWriter, r *http.Request) {
	prIDStr := chi.URLParam(r, "id")
	prID, err := validators.ValidatePullRequestID(prIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	var req dtos.AddReviewerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response_errors.SendError(w, "INVALID_JSON", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateAddReviewerRequest(&req); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	pr, err := h.prService.GetByID(prID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	reviewer, err := h.userService.GetByID(req.ReviewerID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	if err := pr.AddReviewer(reviewer); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	if err := h.prService.Update(pr); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToPullRequestResponse(pr)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *PullRequestHandler) RemoveReviewer(w http.ResponseWriter, r *http.Request) {
	prIDStr := chi.URLParam(r, "id")
	prID, err := validators.ValidatePullRequestID(prIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	reviewerIDStr := chi.URLParam(r, "reviewerID")
	reviewerID, err := validators.ValidateReviewerID(reviewerIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	pr, err := h.prService.GetByID(prID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	if err := pr.RemoveReviewer(reviewerID); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	if err := h.prService.Update(pr); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToPullRequestResponse(pr)
	sendJSONResponse(w, http.StatusOK, response)
}
