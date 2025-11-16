package handlers

import (
	"encoding/json"
	"net/http"
	"reviewer-assignment-service/internal/app/response_errors"
	"reviewer-assignment-service/internal/app/transport/dtos"
	"reviewer-assignment-service/internal/app/transport/mappers"
	"reviewer-assignment-service/internal/app/validators"
	"reviewer-assignment-service/internal/domain/services"

	"github.com/go-chi/chi/v5"
)

type TeamHandler struct {
	teamService services.TeamService
}

func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateTeamRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response_errors.SendError(w, "INVALID_JSON", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateCreateTeamRequest(&req); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}
	team := mappers.ToTeamModel(req)

	if err := h.teamService.Create(team); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToTeamResponse(team)
	sendJSONResponse(w, http.StatusCreated, response)
}

func (h *TeamHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	teamIDStr := chi.URLParam(r, "id")
	teamID, err := validators.ValidateTeamID(teamIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	team, err := h.teamService.GetByID(teamID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToTeamResponse(team)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *TeamHandler) GetTeamByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if err := validators.ValidateTeamName(name); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	team, err := h.teamService.GetByName(name)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToTeamResponse(team)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *TeamHandler) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.teamService.GetAll()
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToTeamListResponse(teams)
	sendJSONResponse(w, http.StatusOK, response)
}

func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	teamIDStr := chi.URLParam(r, "id")
	teamID, err := validators.ValidateTeamID(teamIDStr)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	var req dtos.UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response_errors.SendError(w, "INVALID_JSON", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validators.ValidateUpdateTeamRequest(&req); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	team, err := h.teamService.GetByID(teamID)
	if err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	mappers.UpdateTeamFromRequest(team, req)

	if err := h.teamService.Update(team); err != nil {
		response_errors.HandleServiceError(w, err)
		return
	}

	response := mappers.ToTeamResponse(team)
	sendJSONResponse(w, http.StatusOK, response)
}

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
