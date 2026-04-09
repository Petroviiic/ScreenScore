package storage

import (
	"context"
	"database/sql"
	"time"
)

type DeviceStorage struct {
	db *sql.DB
}

func (d *DeviceStorage) UpdateLastSeen(ctx context.Context, userId int64, deviceId string) error {
	query := `	UPDATE user_devices 
				SET last_seen = NOW() 
				WHERE user_id = $1 AND device_id = $2;`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := d.db.ExecContext(
		ctx,
		query,
		userId,
		deviceId,
	)

	if err != nil {
		return err
	}
	return nil
}

func (d *DeviceStorage) Update(ctx context.Context, userId int64, deviceId string, pushToken string) error {
	query := `
        INSERT INTO user_devices (user_id, device_id, push_token, last_seen)
        VALUES ($1, $2, $3, NOW())
        ON CONFLICT (user_id, device_id) 
        DO UPDATE SET 
            push_token = EXCLUDED.push_token,
            last_seen = NOW();`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := d.db.ExecContext(
		ctx,
		query,
		userId,
		deviceId,
		pushToken,
	)
	return err
}

func (d *DeviceStorage) GetFCMTokens(ctx context.Context, userId int64) ([]string, error) {
	query := `SELECT push_token FROM user_devices WHERE user_id = $1;`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	rows, err := d.db.QueryContext(
		ctx,
		query,
		userId,
	)

	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		token := ""
		err := rows.Scan(
			&token,
		)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (d *DeviceStorage) RequestDeviceSync(ctx context.Context, batchSize int) ([]string, error) {
	query := `
		SELECT DISTINCT ud.push_token
        FROM user_devices AS ud
        WHERE ud.last_seen > NOW() - INTERVAL '7 days' 
        AND ud.push_token IS NOT NULL
        AND ud.push_token != ''
        AND NOT EXISTS (
              SELECT 1 
              FROM screen_time_logs AS stl 
              WHERE stl.user_id = ud.user_id 
                AND stl.recorded_at > NOW() - INTERVAL '1 hour'
          );
	`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := d.db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		token := ""
		err := rows.Scan(
			&token,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (d *DeviceStorage) DeleteFCMToken(token string) error {
	query := `	UPDATE user_devices
				SET push_token = NULL
				WHERE push_token = $1;
			`

	_, err := d.db.Exec(
		query,
		token,
	)
	if err != nil {
		return err
	}
	return nil
}
