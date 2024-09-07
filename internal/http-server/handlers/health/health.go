package health

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/logger/sl"
	resp "url-shortener/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	
}

//go:generate mockery --name=HealthChecker
type HealthChecker interface {
	Health() (bool, error)
}

func New(log *slog.Logger, healthChecker HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.health.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		_, err := healthChecker.Health()
		if err != nil {
			log.Error("Health check failed", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error)

			return
		}

		log.Info("Health check succeded")

		render.Status(r, http.StatusNoContent)
		render.JSON(w, r, resp.OK())
	}
}