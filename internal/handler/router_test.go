package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLHandler(t *testing.T) {
	type args struct {
		urls        map[string]string
		method      string
		targetUrl   string
		contentType string
		body        string
	}
	type want struct {
		needError    bool
		responseCode int
		headerName   string
		headerValue  string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "successful POST request",
			args: args{
				urls:        map[string]string{},
				method:      http.MethodPost,
				targetUrl:   "/",
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
				targetUrl:   "/wrong",
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
				targetUrl:   "/",
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
				targetUrl: "/sULftRJq",
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
				targetUrl: "/sULftRJq/wrong",
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
				targetUrl: "/sULftRJq",
			},
			want: want{
				needError:    true,
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name: "unsuccessfu request, wrong http method",
			args: args{
				method:    http.MethodPut,
				targetUrl: "/sULftRJq",
			},
			want: want{
				needError:    true,
				responseCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(tt.args.body)
			request := httptest.NewRequest(tt.args.method, tt.args.targetUrl, bytes.NewBuffer(body))
			request.Header.Set("Content-Type", tt.args.contentType)
			w := httptest.NewRecorder()
			URLHandler(tt.args.urls)(w, request)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.want.responseCode, resp.StatusCode)

			if !tt.want.needError {
				if tt.args.method == http.MethodPost {
					resBody, err := io.ReadAll(resp.Body)
					require.NoError(t, err)
					resParts := strings.Split(string(resBody), "/")
					_, ok := tt.args.urls[resParts[len(resParts)-1]]
					assert.True(t, ok)
				}
				if tt.args.method == http.MethodGet {
					headerValue := resp.Header.Get(tt.want.headerName)
					assert.Equal(t, tt.want.headerValue, headerValue)
				}
			}
		})
	}
}
