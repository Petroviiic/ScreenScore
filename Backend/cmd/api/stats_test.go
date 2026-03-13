package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/Petroviiic/ScreenScore/internal/storage"
)

func TestValidateScreenTime(t *testing.T) {
	// app := &Application{
	// 	config:  Config{},
	// 	storage: storage.NewMockStorage(),
	// }
	// mux := app.mount()

	type testCase struct {
		name         string
		lastRecord   storage.UsageRecord
		currentStats storage.UsageRecord
		currentTime  time.Time
		wantErr      bool
	}

	tests := []testCase{
		{
			name:         "Validan rast unutar istog dana",
			lastRecord:   storage.UsageRecord{ScreenTime: 100, RecordedAt: time.Now().UTC().Add(-1 * time.Hour)},
			currentStats: storage.UsageRecord{ScreenTime: 130}, // 30 min rasta u 60 min vremena
			currentTime:  time.Now().UTC(),
			wantErr:      false,
		},
		{
			name:         "Prevara: Brži rast od realnog vremena",
			lastRecord:   storage.UsageRecord{ScreenTime: 100, RecordedAt: time.Now().UTC().Add(-10 * time.Minute)},
			currentStats: storage.UsageRecord{ScreenTime: 120}, // 20 min rasta u 10 min vremena - NEMOGUĆE
			currentTime:  time.Now().UTC(),
			wantErr:      true,
		},
		{
			name:         "Reset: Novi dan (manje minuta nego juče)",
			lastRecord:   storage.UsageRecord{ScreenTime: 500, RecordedAt: time.Now().UTC().Add(-24 * time.Hour)},
			currentStats: storage.UsageRecord{ScreenTime: 10}, // Novi dan, krenuo od nule
			currentTime:  time.Now().UTC(),
			wantErr:      false,
		},
		{
			name:         "Budućnost: Sat na telefonu pomjeren unaprijed",
			lastRecord:   storage.UsageRecord{ScreenTime: 100, RecordedAt: time.Now().UTC()},
			currentStats: storage.UsageRecord{ScreenTime: 110},
			currentTime:  time.Now().UTC().Add(2 * time.Hour), // 2 sata u budućnosti
			wantErr:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			http.NewRequest("POST", "/v1/stats/sync-stats", nil)
		})
	}
}
