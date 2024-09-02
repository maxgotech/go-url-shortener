package delete_test

import (
	"encoding/json"
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
		code      int
	}{
		{
			name:  "Success",
			alias: "alias",
			code:  http.StatusOK,
		},
		{
			name:      "Alias Not Found",
			alias:     "abcde",
			respError: "Alias doesnt exist",
			mockError: storage.ErrURLNotFound,
			code:      http.StatusNotFound,
		},
		{
			name:      "Delete Error",
			alias:     "alias",
			respError: "failed to delete Alias",
			mockError: errors.New("unexpected error"),
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias, mock.AnythingOfType("string")).
					Return(true, tc.mockError).
					Once()
			}

			handler := chi.NewRouter()
			handler.Delete("/{alias}", delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			req, err := http.NewRequest(http.MethodDelete, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.code)

			body := rr.Body.String()

			var resp resp.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
