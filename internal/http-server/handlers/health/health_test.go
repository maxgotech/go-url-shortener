package health_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/health"
	"url-shortener/internal/http-server/handlers/health/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/stretchr/testify/require"
)

func TestHealthChecker(t *testing.T) {
	cases := []struct {
		name      string
		respError string
		mockError error
		code      int
	}{
		{
			name: "ok",
			code: http.StatusNoContent,
		},
		{
			name:      "error",
			respError: "Health check failed",
			mockError: errors.New("unexpected error"),
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			healthCheckerMock := mocks.NewHealthChecker(t)

			if tc.respError == "" || tc.mockError != nil {
				healthCheckerMock.On("Health").
					Return(true, tc.mockError).Once()
			}

			handler := health.New(slogdiscard.NewDiscardLogger(), healthCheckerMock)

			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.code)

		})
	}

}
