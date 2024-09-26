package save

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate mockery --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

type KafkaProducer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

func New(log *slog.Logger, prod KafkaProducer, urlSaver URLSaver, alias_length int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		// TODO(Maxim): check if alias already exists
		alias := req.Alias
		if alias == "" {
			alias = random.RandStr(alias_length)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("URL already exists", slog.String("url", req.URL))

			render.Status(r, http.StatusConflict)
			render.JSON(w, r, resp.Error("URL already exists"))

			return
		}

		if err != nil {
			log.Error("failed to save URL", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to save URL"))

			return
		}

		log.Info("URL added", slog.Int64("id", id))

		msg := kafka.Message{
			Key:   []byte(strconv.Itoa(int(id))),
			Value: []byte(fmt.Sprintf("%v with alias %v", req.URL, req.Alias)),
		}

		err = prod.WriteMessages(context.Background(), msg)
		if err != nil {
			log.Warn("failed to log to kafka", sl.Err(err))
		} else {
			log.Info("Kafka log send")
		}

		render.Status(r, http.StatusCreated)
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
