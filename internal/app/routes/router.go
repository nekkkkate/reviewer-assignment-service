package routes

import (
	"net/http"
	"reviewer-assignment-service/internal/app/handlers"

	"reviewer-assignment-service/internal/domain/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(
	userService services.UserService,
	prService services.PullRequestService,
	teamService services.TeamService,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	userHandler := handlers.NewUserHandler(userService, prService)
	teamHandler := handlers.NewTeamHandler(teamService)

	r.Route("/", func(r chi.Router) {

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Post("/setIsActive", userHandler.SetUserActive)
			r.Post("/deactivate", userHandler.DeactivateUser)
			r.Get("/getReview", userHandler.GetUserReviewPRs)
			r.Get("/by-email", userHandler.GetUserByEmail)
			r.Get("/", userHandler.GetAllUsers)
			r.Get("/{id}", userHandler.GetUserByID)
		})

		r.Route("/teams", func(r chi.Router) {
			r.Get("/", teamHandler.GetAllTeams)
			r.Post("/", teamHandler.CreateTeam)
			r.Get("/by-name/{name}", teamHandler.GetTeamByName)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", teamHandler.GetTeamByID)
				r.Put("/", teamHandler.UpdateTeam)
			})
		})

	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	return r
}
