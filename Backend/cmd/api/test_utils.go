package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Petroviiic/ScreenScore/internal/auth"
	"github.com/Petroviiic/ScreenScore/internal/ratelimiter"
	"github.com/Petroviiic/ScreenScore/internal/storage"
)

func newTestApplication(t *testing.T) *Application {
	t.Helper()

	storage := storage.NewMockStorage()
	auth := auth.NewMockJWTAuthenticator()
	return &Application{
		storage:       storage,
		authenticator: auth,
		rateLimiters: rateLimiters{
			apiFixedWindow:  ratelimiter.NewFixedWindowLimiter(100, 1),
			authFixedWindow: ratelimiter.NewFixedWindowLimiter(100, 1),
			tokenBucket:     ratelimiter.NewTokenBuckerRatelimiter(10, 10),
		},
		config: Config{
			userStreak: userStreakConfig{
				maxShieldCount:          5,
				minScreenTimeThreshold:  120,
				shieldCountIncreaseRate: 5,
			},
		},
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected the response code to be %d and we got %d", expected, actual)
	}
}
