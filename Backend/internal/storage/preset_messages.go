package storage

import (
	"context"
	"database/sql"
)

type PresetMessageStorage struct {
	db *sql.DB
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
