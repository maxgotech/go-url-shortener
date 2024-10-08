package delete

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate mockery --name=URLDeleter
type URLDeleter interface {
	DeleteURL(urlToDelete string) (bool, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		// should never shoot(atleast with chi) since
		// if path is empty router throws 404 page not found
		// without calling the func
		// but just in case this will stay
		if alias == "" {
			log.Info("alias is empty")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		log.Info("got alias from url", slog.String("Alias", alias))

		_, err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("URL with that Alias doesnt exist", slog.String("Alias", alias))

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("Alias doesnt exist"))

			return
		}

		if err != nil {
			log.Error("failed to delete Alias", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete Alias"))

			return
		}

		log.Info("Alias deleted", slog.String("Alias", alias))

		render.Status(r, http.StatusOK)
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
