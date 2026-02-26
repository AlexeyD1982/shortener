package handler

import (
	"net/http"
	"reflect"
	"testing"
)

func TestURLHandler(t *testing.T) {
	type args struct {
		urls        map[string]string
		method      string
		target      string
		contentType string
	}
	tests := []struct {
		name string
		args args
		want func(w http.ResponseWriter, r *http.Request)
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLHandler(tt.args.urls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URLHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
