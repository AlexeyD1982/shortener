package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generateRandomString(t *testing.T) {
	type args struct {
		length int
		seed   int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "short seed",
			args: args{
				length: 8,
				seed:   999,
			},
			want: "qUtjrsDU",
		},
		{
			name: "long seed",
			args: args{
				length: 218,
				seed:   999654564654465,
			},
			want: "vdUWmUVznyWNYrAhPbMpzExdWvzBDYYvXBRhqeRgYYqguLKsBxotPYMGbXuIMRSoCiaHubTCrldFXVCWRQwgjBMoyuyRiwDGeucoJeZorlEHbjaWRilwOmluQYOaGGRDaZGgHgiYygdPyQHlGUSnsgoaQtFnPFzSHSHQvmzWeyGiQZXtcWEcQYTiLfbWMrJQITErVLYTxQgbmBbahgBjvQgxLb",
		},
		{
			name: "long length",
			args: args{
				length: 8,
				seed:   999654564654465,
			},
			want: "vdUWmUVz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, generateRandomString(tt.args.length, tt.args.seed), "generateRandomString(%v, %v)", tt.args.length, tt.args.seed)
		})
	}
}
