package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexeyD1982/shortener/internal/config"
	"github.com/AlexeyD1982/shortener/internal/repository/inmemory"
	"github.com/AlexeyD1982/shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path, contentType, body string) (*http.Response, string) {
	t.Helper()
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer([]byte(body)))
	require.NoError(t, err)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestURLRouter(t *testing.T) {
	type args struct {
		urls        map[string]string
		method      string
		targetURL   string
		contentType string
		body        string
	}
	type want struct {
		needError    bool
		responseCode int
		headerName   string
		headerValue  string
	}
	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "successful POST request",
			args: args{
				urls:        map[string]string{},
				method:      http.MethodPost,
				targetURL:   "/",
				contentType: "text/plain",
				body:        "https://test.ru",
			},
			want: want{
				needError:    false,
				responseCode: http.StatusCreated,
			},
		},
		{
			name: "unsuccessful POST request, wrong path",
			args: args{
				urls:        map[string]string{},
				method:      http.MethodPost,
				targetURL:   "/wrong",
				contentType: "text/plain",
				body:        "https://test.ru",
			},
			want: want{
				needError:    true,
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name: "unsuccessful POST request, wrong content type",
			args: args{
				urls:        map[string]string{},
				method:      http.MethodPost,
				targetURL:   "/",
				contentType: "application/json",
				body:        "https://test.ru",
			},
			want: want{
				needError:    true,
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name: "successful GET request",
			args: args{
				urls:      map[string]string{"sULftRJq": "https://test.ru"},
				method:    http.MethodGet,
				targetURL: "/sULftRJq",
			},
			want: want{
				needError:    false,
				responseCode: http.StatusTemporaryRedirect,
				headerName:   "Location",
				headerValue:  "https://test.ru",
			},
		},
		{
			name: "unsuccessful GET request, wrong path",
			args: args{
				urls:      map[string]string{"sULftRJq": "https://test.ru"},
				method:    http.MethodGet,
				targetURL: "/sULftRJq/wrong",
			},
			want: want{
				needError:    true,
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name: "unsuccessful GET request, no saved url",
			args: args{
				urls:      map[string]string{},
				method:    http.MethodGet,
				targetURL: "/sULftRJq",
			},
			want: want{
				needError:    true,
				responseCode: http.StatusNotFound,
			},
		},
		{
			name: "unsuccessfu request, wrong http method",
			args: args{
				urls:      map[string]string{},
				method:    http.MethodPut,
				targetURL: "/sULftRJq",
			},
			want: want{
				needError:    true,
				responseCode: http.StatusBadRequest,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Get(false)

			repo := inmemory.NewStorageWithData(tc.args.urls)
			srv := service.NewURLService(repo)

			core, _ := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			r := NewURLHandler(srv, cfg, logger)
			ts := httptest.NewServer(r.InitRouter())
			defer ts.Close()
			resp, body := testRequest(t, ts, tc.args.method, tc.args.targetURL, tc.args.contentType, tc.args.body)
			defer resp.Body.Close()
			assert.Equal(t, tc.want.responseCode, resp.StatusCode)

			if !tc.want.needError {
				if tc.args.method == http.MethodPost {
					resParts := strings.Split(body, "/")
					_, ok := tc.args.urls[resParts[len(resParts)-1]]
					assert.True(t, ok)
				}
				if tc.args.method == http.MethodGet {
					headerValue := resp.Header.Get(tc.want.headerName)
					assert.Equal(t, tc.want.headerValue, headerValue)
				}
			}
		})
	}
}
