package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/kafka/producerdiscard"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {

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
			url:   "https://google.com",
			code:  http.StatusCreated,
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
			code:  http.StatusCreated,
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
			code:      http.StatusBadRequest,
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid url",
			code:      http.StatusBadRequest,
		},
		{
			name:      "URL exists",
			alias:     "existing_alias",
			url:       "https://google.com",
			respError: "URL already exists",
			mockError: storage.ErrURLExists,
			code:      http.StatusConflict,
		},
		{
			name:      "SaveURL Error",
			alias:     "error_alias",
			url:       "https://google.com",
			respError: "failed to save URL",
			mockError: errors.New("unexpected error"),
			code:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			// _, _, 5 is alias_length
			handler := save.New(slogdiscard.NewDiscardLogger(), producerdiscard.NewDiscardProducer(), urlSaverMock, 5)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.code)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

			// TODO: add more checks
		})
	}

}
