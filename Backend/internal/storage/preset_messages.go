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
	Message   string    `json:"message"`
	Price     int       `json:"price"`
	Rarity    string    `json:"rarity"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *PresetMessageStorage) InsertNewPresetMessage(ctx context.Context, tx *sql.Tx, text string, price int, rarity string, isActive bool) error {
	return func() error {
		query := `
					INSERT INTO preset_messages (message, price, rarity, is_active)
					VALUES ($1, $2, $3, $4)
					ON CONFLICT (message) 
					DO UPDATE SET 
						price = EXCLUDED.price,
						rarity = EXCLUDED.rarity,
						is_active = EXCLUDED.is_active;;
				`

		_, err := tx.ExecContext(
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
	}()
}

func (m *PresetMessageStorage) GetPresetMessage(ctx context.Context, userID int64, msgID int64) (string, error) {
	query := `
			SELECT pm.message from user_messages AS um 
			JOIN preset_messages AS pm ON pm.id = um.message_id
			WHERE um.user_id = $1 AND um.message_id = $2;
		`

	msg := ""
	err := m.db.QueryRowContext(
		ctx,
		query,
		userID,
		msgID,
	).Scan(
		&msg,
	)

	if err != nil {
		return "", err
	}

	return msg, nil
}

func (m *PresetMessageStorage) GetOwnedPresetMessage(ctx context.Context, userID int64) ([]*PresetMessage, error) {
	query := `		
			SELECT pm.id, pm.message, pm.price, pm.rarity, pm.is_active, pm.created_at from user_messages AS um 
			JOIN preset_messages AS pm ON pm.id = um.message_id
			WHERE um.user_id = $1;
		`
	rows, err := m.db.QueryContext(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var msgs []*PresetMessage
	for rows.Next() {
		msg := &PresetMessage{}
		err := rows.Scan(
			&msg.ID,
			&msg.Message,
			&msg.Price,
			&msg.Rarity,
			&msg.IsActive,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (m *PresetMessageStorage) GetAvaiableInShop(ctx context.Context, userID int64) ([]*PresetMessage, error) {
	query := `
			SELECT pm.id, pm.message, pm.price, pm.rarity, pm.is_active, pm.created_at FROM preset_messages AS pm 
			LEFT JOIN user_messages AS um ON pm.id = um.message_id AND um.user_id = $1
			WHERE um.user_id IS NULL AND pm.is_active = TRUE;
		`

	rows, err := m.db.QueryContext(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*PresetMessage
	for rows.Next() {
		msg := &PresetMessage{}
		err := rows.Scan(
			&msg.ID,
			&msg.Message,
			&msg.Price,
			&msg.Rarity,
			&msg.IsActive,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}
