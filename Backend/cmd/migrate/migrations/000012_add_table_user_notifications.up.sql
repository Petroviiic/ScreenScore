CREATE TABLE IF NOT EXISTS user_notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    points_earned NUMERIC(10, 2) DEFAULT 0.00,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    notification_type VARCHAR(50),                                  -- eg. 'weekly_reward', 'added_to_group'
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_notifications_unread ON user_notifications(user_id) WHERE is_read = FALSE;