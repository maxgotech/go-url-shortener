package redirect_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/redirect/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
		code      int
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.github.com",
			code:  http.StatusFound,
		},
		{
			name:      "Error not found",
			alias:     "test",
			url:       "https://www.github.com",
			respError: "not found",
			mockError: storage.ErrURLNotFound,
			code:      http.StatusNotFound,
		},
		{
			name:      "unexpected error",
			alias:     "alias",
			respError: "failed to delete Alias",
			mockError: errors.New("unexpected error"),
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			handler := chi.NewRouter()
			handler.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			req, err := http.NewRequest(http.MethodGet, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.code)

			if tc.code == http.StatusFound {
				require.Equal(t, rr.Header().Get("Location"), tc.url)
			}
		})
	}
}
