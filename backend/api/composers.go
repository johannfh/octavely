package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/johannfh/octavely/backend/httputils"
)

func (s *Server) handleGetComposer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	idParam := chi.URLParam(r, "id")

	if idParam == "" {
		res := httputils.NewResponse(
			httputils.WithStatusCode[any](http.StatusBadRequest),
			httputils.WithResponseError[any](httputils.NewResponseError(
				httputils.WithErrorCode("INVALID_ID_PARAM"),
				httputils.WithErrorMessage("parameter 'id' is invalid"),
				httputils.WithErrorDetails("parameter 'id' is empty but must be a valid integer"),
				httputils.WithErrorTimestamp(time.Now()),
				httputils.WithErrorPath(r.URL.Path),
			)),
		)
		httputils.WriteJson(w, http.StatusBadRequest, res)
		return
	}

	id, err := strconv.Atoi(idParam)

	if err != nil {
		res := httputils.NewResponse(
			httputils.WithStatusCode[any](http.StatusBadRequest),
			httputils.WithResponseError[any](httputils.NewResponseError(
				httputils.WithErrorCode("INVALID_ID_PARAM"),
				httputils.WithErrorMessage("parameter 'id' is invalid"),
				httputils.WithErrorDetails(fmt.Sprintf("parameter 'id' is '%s' but must be a valid integer integer", idParam)),
				httputils.WithErrorTimestamp(time.Now()),
				httputils.WithErrorPath(r.URL.Path),
			)),
		)
		httputils.WriteJson(w, http.StatusBadRequest, res)
		return
	}

	composer, err := s.queries.GetComposer(ctx, int64(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res := httputils.NewResponse(
				httputils.WithStatusCode[any](http.StatusNotFound),
				httputils.WithResponseError[any](httputils.NewResponseError(
					httputils.WithErrorCode("COMPOSER_NOT_FOUND"),
					httputils.WithErrorMessage("composer not found"),
					httputils.WithErrorDetails(fmt.Sprintf("composer with the ID '%d' was not found", id)),
					httputils.WithErrorTimestamp(time.Now()),
					httputils.WithErrorPath(r.URL.Path),
				)),
			)
			httputils.WriteJson(w, http.StatusNotFound, res)
			return
		}

		s.logger.Error("failed to retrieve composer from database", "err", err, "composerId", id)

		res := httputils.NewResponse(
			httputils.WithStatusCode[any](http.StatusInternalServerError),
			httputils.WithResponseError[any](httputils.NewResponseError(
				httputils.WithErrorCode("INTERNAL_SERVER_ERROR"),
				httputils.WithErrorMessage("internal server error"),
				httputils.WithErrorTimestamp(time.Now()),
				httputils.WithErrorPath(r.URL.Path),
			)),
		)
		httputils.WriteJson(w, http.StatusInternalServerError, res)
		return
	}

	type composerData struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}

	res := httputils.NewResponse(
		httputils.WithStatusCode[composerData](http.StatusOK),
		httputils.WithResponseData(&composerData{
			Id:   composer.ID,
			Name: composer.Name,
		}),
	)

	httputils.WriteJson(w, http.StatusOK, res)
	slog.Info(
		fmt.Sprintf("%s %s", r.Method, r.URL.Path),
		"method", r.Method,
		"path", r.URL.Path,
		"response", res,
	)
}

func (s *Server) handleGetAllComposers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	composers, err := s.queries.GetAllComposers(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		httputils.WriteJson(w, http.StatusInternalServerError, struct{}{})
		return
	}

	type composerData struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}

	data := []composerData{}
	for _, composer := range composers {
		data = append(data, composerData{
			Id:   composer.ID,
			Name: composer.Name,
		})
	}

	res := httputils.NewResponse(
		httputils.WithStatusCode[[]composerData](http.StatusOK),
		httputils.WithResponseData(&data),
	)
	httputils.WriteJson(w, http.StatusOK, res)

	slog.Info(
		fmt.Sprintf("%s %s", r.Method, r.URL.Path),
		"method", r.Method,
		"path", r.URL.Path,
		"response", res,
	)
}
