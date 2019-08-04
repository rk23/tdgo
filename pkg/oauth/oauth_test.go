package oauth

import (
	"reflect"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	type args struct {
		req *AccessTokenRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *AccessTokenResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RefreshToken(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
