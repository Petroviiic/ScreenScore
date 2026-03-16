package storage

import (
	"context"
	"database/sql"
)

type DeviceStorage struct {
	db *sql.DB
}

func (d *DeviceStorage) Update(ctx context.Context, userId int64, deviceId string, pushToken string) error {
	query := `
        INSERT INTO user_devices (user_id, device_id, push_token, last_seen)
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (user_id, device_id) 
        DO UPDATE SET 
            push_token = EXCLUDED.push_token,
            last_seen = NOW();`

	_, err := d.db.ExecContext(
		ctx,
		query,
		userId,
		deviceId,
		pushToken,
	)
	return err
}
