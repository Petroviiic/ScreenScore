CREATE TABLE IF NOT EXISTS screen_time_logs(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    screen_time INT NOT NULL,
    recorded_at TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_screen_time_user_date ON screen_time_logs(user_id, recorded_at DESC);