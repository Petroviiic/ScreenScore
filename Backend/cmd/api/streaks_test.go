package main

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Petroviiic/ScreenScore/internal/storage"
)

func TestSyncStreak(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)
	threeDaysAgo := now.AddDate(0, 0, -3)

	currYear, currWeek := now.ISOWeek()
	lastYear, lastWeek := now.AddDate(0, 0, -7).ISOWeek()

	type testCase struct {
		name                     string
		averageScreenTimeForWeek float64
		screenTimeForDay         int
		shieldCount              int
		lastRecordedWeekNumber   int
		lastRecordedYearNumber   int
		lastWeekAverage          float64
		lastUpdatedAt            time.Time
		expectedStatusCode       int
		expectedFrozen           bool
		expectedShieldsNeeded    int
		nowTime                  time.Time
	}

	tests := []testCase{
		{
			name:                     "Success: Good screentime yesterday, streak incremented",
			averageScreenTimeForWeek: 120,
			screenTimeForDay:         60,
			shieldCount:              1,
			lastRecordedWeekNumber:   currWeek,
			lastRecordedYearNumber:   currYear,
			lastWeekAverage:          120,
			lastUpdatedAt:            yesterday,
			expectedStatusCode:       http.StatusOK,
			expectedFrozen:           false,
			expectedShieldsNeeded:    0,
			nowTime:                  time.Now().UTC(),
		},
		{
			name:                     "Frozen: Bad screentime yesterday, needs 1 shield",
			averageScreenTimeForWeek: 120,
			screenTimeForDay:         180,
			shieldCount:              2,
			lastRecordedWeekNumber:   currWeek,
			lastRecordedYearNumber:   currYear,
			lastWeekAverage:          120,
			lastUpdatedAt:            yesterday,
			expectedStatusCode:       http.StatusOK,
			expectedFrozen:           true,
			expectedShieldsNeeded:    1,
			nowTime:                  time.Now().UTC(),
		},
		{
			name:                     "Frozen: Inactivity gap (3 days), needs 2 shields",
			averageScreenTimeForWeek: 120,
			screenTimeForDay:         60,
			shieldCount:              5,
			lastRecordedWeekNumber:   currWeek,
			lastRecordedYearNumber:   currYear,
			lastWeekAverage:          120,
			lastUpdatedAt:            threeDaysAgo,
			expectedStatusCode:       http.StatusOK,
			expectedFrozen:           true,
			expectedShieldsNeeded:    2,
			nowTime:                  time.Now().UTC(),
		},
		{
			name:                     "Week Transition: Updates average from last week",
			averageScreenTimeForWeek: 300,
			screenTimeForDay:         200,
			shieldCount:              1,
			lastRecordedWeekNumber:   lastWeek - 1,
			lastRecordedYearNumber:   lastYear,
			lastWeekAverage:          100,
			lastUpdatedAt:            yesterday,
			expectedStatusCode:       http.StatusOK,
			expectedFrozen:           false,
			expectedShieldsNeeded:    0,
			nowTime:                  time.Now().UTC(),
		},
		{
			name:                     "Today already processed",
			averageScreenTimeForWeek: 0,
			screenTimeForDay:         0,
			shieldCount:              0,
			lastRecordedWeekNumber:   0,
			lastRecordedYearNumber:   0,
			lastWeekAverage:          0,
			lastUpdatedAt:            now,
			expectedStatusCode:       http.StatusOK,
			expectedFrozen:           false,
			expectedShieldsNeeded:    0,
			nowTime:                  time.Now().UTC(),
		},

		{
			name:                     "New Year Transition: Same ISO Week (Dec 31 to Jan 1)",
			averageScreenTimeForWeek: 150,
			screenTimeForDay:         100,
			shieldCount:              1,
			lastRecordedWeekNumber:   1,
			lastRecordedYearNumber:   2026,
			lastWeekAverage:          150,
			lastUpdatedAt:            time.Date(2025, 12, 31, 23, 0, 0, 0, time.UTC),
			expectedStatusCode:       http.StatusOK,
			expectedFrozen:           false,
			expectedShieldsNeeded:    0,
			nowTime:                  time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		{
			name:                     "New Year: First Monday of January (Update Average)",
			averageScreenTimeForWeek: 250,
			screenTimeForDay:         100,
			shieldCount:              3,
			lastRecordedWeekNumber:   52,
			lastRecordedYearNumber:   2025,
			lastWeekAverage:          200,
			lastUpdatedAt:            time.Date(2026, 1, 4, 20, 0, 0, 0, time.UTC),
			expectedStatusCode:       http.StatusOK,
			expectedFrozen:           false,
			expectedShieldsNeeded:    0,
			nowTime:                  time.Date(2026, 1, 5, 10, 0, 0, 0, time.UTC),
		},
		{
			name:                     "New Year: Long Gap Over Transition (Dec 30 to Jan 2)",
			averageScreenTimeForWeek: 150,
			screenTimeForDay:         100,
			shieldCount:              5,
			lastRecordedWeekNumber:   1,
			lastRecordedYearNumber:   2026,
			lastWeekAverage:          150,
			lastUpdatedAt:            time.Date(2025, 12, 30, 10, 0, 0, 0, time.UTC),
			expectedStatusCode:       http.StatusOK,
			expectedFrozen:           true,
			expectedShieldsNeeded:    2,
			nowTime:                  time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStats := app.storage.StatsStorage.(*storage.StatsMockStorage)

			mockStats.GetUserAverageScreenTimeForWeekFunc = func(ctx context.Context, userID int64) (float64, error) {
				return tc.averageScreenTimeForWeek, nil
			}
			mockStats.GetUserScreenTimeForDayFunc = func(ctx context.Context, userID int64) (int, error) {
				return tc.screenTimeForDay, nil
			}

			mockUserStreaks := app.storage.UserStreakStorage.(*storage.UserStreakMockStorage)

			mockUserStreaks.GetStreakDataFunc = func(ctx context.Context, userID int64) (*storage.StreakData, error) {
				return &storage.StreakData{
					CurrentStreak:   10,
					AllTimeHigh:     20,
					ShieldCount:     tc.shieldCount,
					WeekNumber:      tc.lastRecordedWeekNumber,
					YearNumber:      tc.lastRecordedYearNumber,
					LastWeekAverage: tc.lastWeekAverage,
					LastUpdatedAt:   tc.lastUpdatedAt,
				}, nil
			}

			app.clock.(*MockClock).FixedTime = tc.nowTime
			req, err := http.NewRequest(http.MethodPost, "/v1/streak/sync", nil)
			req.Header.Add("Authorization", "Bearer 123token")
			if err != nil {
				t.Fatal(err)
			}

			resp := executeRequest(req, mux)

			checkResponseCode(t, tc.expectedStatusCode, resp.Code)

			var result struct {
				Data struct {
					StreakValidation *StreakValidation              `json:"streak_validation"`
					StreakData       *storage.UserStreakMockStorage `json:"streak_data"`
				} `json:"data"`
			}

			_ = json.NewDecoder(resp.Body).Decode(&result)
			if result.Data.StreakValidation.StreakFrozen != tc.expectedFrozen {
				t.Errorf("expected frozen %v, got %v", tc.expectedFrozen, result.Data.StreakValidation.StreakFrozen)
			}

			if result.Data.StreakValidation.ShieldsNeeded != tc.expectedShieldsNeeded {
				t.Errorf("expected shields needed %d, got %d", tc.expectedShieldsNeeded, result.Data.StreakValidation.ShieldsNeeded)
			}
		})
	}
}
