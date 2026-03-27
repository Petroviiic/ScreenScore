package storage

import (
	"context"
	"database/sql"
	"time"
)

type PresetMessageStorage struct {
	db *sql.DB
}
type PresetMessage struct {
	ID        int64     `json:"id"`
	Price     int       `json:"price"`
	Rarity    string    `json:"rarity"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *PresetMessageStorage) InsertNewPresetMessage(ctx context.Context, text string, price int, rarity string, isActive bool) error {
	query := `
				INSERT INTO preset_messages (message, price, rarity, is_active)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (message) 
				DO UPDATE SET 
					price = EXCLUDED.price,
					rarity = EXCLUDED.rarity,
					is_active = EXCLUDED.is_active;;
			`

	_, err := m.db.ExecContext(
		ctx,
		query,
		text,
		price,
		rarity,
		isActive,
	)

	if err != nil {
		return err
	}
	return nil
}
