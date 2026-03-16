package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Petroviiic/ScreenScore/internal/storage"
)

func TestValidateScreenTime(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	type testCase struct {
		name               string
		lastRecord         storage.UsageRecord
		currentScreenTime  int32
		currentTime        string
		expectedStatusCode int
		deviceId           string
	}

	refTime, _ := time.Parse(time.RFC3339, "2026-03-13T12:00:00+01:00")
	refTime = refTime.UTC()

	tests := []testCase{
		{
			name:               "30 mins of screentime in 1 hour. Ok",
			lastRecord:         storage.UsageRecord{ScreenTime: 100, RecordedAt: refTime.Add(-1 * time.Hour)},
			currentScreenTime:  130,
			currentTime:        refTime.Format(time.RFC3339),
			expectedStatusCode: http.StatusCreated,
			deviceId:           "testdevice",
		},
		{
			name:               "20 mins of screentime in 10 mins. Impossible. Deny request",
			lastRecord:         storage.UsageRecord{ScreenTime: 100, RecordedAt: refTime.Add(-10 * time.Minute)},
			currentScreenTime:  120,
			currentTime:        refTime.Format(time.RFC3339),
			expectedStatusCode: http.StatusBadRequest,
			deviceId:           "testdevice",
		},
		{
			name:               "New record sent before the last saved. Deny request",
			lastRecord:         storage.UsageRecord{ScreenTime: 100, RecordedAt: refTime.Add(10 * time.Minute)},
			currentScreenTime:  120,
			currentTime:        refTime.Format(time.RFC3339),
			expectedStatusCode: http.StatusBadRequest,
			deviceId:           "testdevice",
		},
		{
			name:               "New screentime longer than duration of the current day. Deny",
			lastRecord:         storage.UsageRecord{ScreenTime: 150, RecordedAt: refTime},
			currentScreenTime:  800,
			currentTime:        refTime.Format(time.RFC3339),
			expectedStatusCode: http.StatusBadRequest,
			deviceId:           "testdevice",
		},
		{
			name:               "Lower screen time than the last record the same day. Deny",
			lastRecord:         storage.UsageRecord{ScreenTime: 150, RecordedAt: refTime.Add(-1 * time.Hour)},
			currentScreenTime:  130,
			currentTime:        refTime.Format(time.RFC3339),
			expectedStatusCode: http.StatusBadRequest,
			deviceId:           "testdevice",
		},
		{
			name:               "Reset: Accept lower screen time, because it's a different day",
			lastRecord:         storage.UsageRecord{ScreenTime: 500, RecordedAt: refTime.Add(-24 * time.Hour)},
			currentScreenTime:  10,
			currentTime:        refTime.Format(time.RFC3339),
			expectedStatusCode: http.StatusCreated,
			deviceId:           "testdevice",
		},
		{
			name:               "Record sent from the future. Deny",
			lastRecord:         storage.UsageRecord{ScreenTime: 100, RecordedAt: time.Now().UTC()},
			currentScreenTime:  110,
			currentTime:        time.Now().UTC().Add(2 * time.Hour).Format(time.RFC3339),
			expectedStatusCode: http.StatusBadRequest,
			deviceId:           "testdevice",
		},
		{
			name:               "Screen time is the same as the last one. Deny",
			lastRecord:         storage.UsageRecord{ScreenTime: 100, RecordedAt: refTime.Add(-1 * time.Hour)},
			currentScreenTime:  100,
			currentTime:        refTime.Format(time.RFC3339),
			expectedStatusCode: http.StatusBadRequest,
			deviceId:           "testdevice",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := app.storage.StatsStorage.(*storage.StatsMockStorage)

			mock.GetUsersLastFunc = func(ctx context.Context, userID int64) (*storage.UsageRecord, error) {
				return &storage.UsageRecord{
					ScreenTime: tc.lastRecord.ScreenTime,
					RecordedAt: tc.lastRecord.RecordedAt,
					DeviceID:   tc.deviceId,
				}, nil
			}
			jsonBody := fmt.Sprintf(`{"screen_time": %d, "recorded_at": "%s", "device_id":"%s"}`, tc.currentScreenTime, tc.currentTime, tc.deviceId)
			req, err := http.NewRequest(http.MethodPost, "/v1/stats/sync-stats", strings.NewReader(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			resp := executeRequest(req, mux)

			checkResponseCode(t, tc.expectedStatusCode, resp.Code)
		})
	}
}
