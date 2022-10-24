package token_test

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/achristie/save2db/internal/assert"
	"github.com/achristie/save2db/pkg/platts/token"
)

var (
	//go:embed testdata/429.json
	RespRateLimit string
	//go:embed testdata/400.json
	RespInvalidCred string
	//go:embed testdata/401.json
	RespMissingAPIKey string
	//go:embed testdata/200.json
	RespOK string
)

func TestGetToken(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    error
	}{
		{
			name:       "Rate Limit Exceeded",
			statusCode: http.StatusTooManyRequests,
			wantErr:    token.ErrRateLimited,
			body:       RespRateLimit,
		},
		{
			name:       "Invalid Credentials",
			statusCode: http.StatusBadRequest,
			wantErr:    token.ErrInvalidCred,
			body:       RespInvalidCred,
		},
		{
			name:       "Missing API Key",
			statusCode: http.StatusUnauthorized,
			wantErr:    token.ErrInvalidCred,
			body:       RespMissingAPIKey,
		},
		{
			name:       "Server Error",
			statusCode: http.StatusBadGateway,
			wantErr:    token.ErrServerIssue,
		},
		{
			name:       "OK",
			statusCode: http.StatusOK,
			body:       RespOK,
		},
		{
			name:       "Cache",
			statusCode: http.StatusOK,
			body:       RespOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.body))
			}))
			defer ts.Close()

			tc := token.NewTokenClient("test", "test", "test")
			tc.TokenEndpoint = ts.URL

			res, err := tc.GetToken()
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			var tkn token.Token
			json.Unmarshal([]byte(tt.body), &tkn)

			assert.Equal(t, res.AccessToken, tkn.AccessToken)

			//make sure the IAT matches for subsequent calls
			time.Sleep(10 * time.Millisecond)
			resCache, _ := tc.GetToken()
			assert.Equal(t, resCache.Iat, res.Iat)
		})
	}
}
