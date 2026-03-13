package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Petroviiic/ScreenScore/internal/storage"
)

func newTestApplication(t *testing.T) *Application {
	t.Helper()

	storage := storage.NewMockStorage()

	return &Application{
		storage: storage,
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
