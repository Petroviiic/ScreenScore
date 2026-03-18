DROP INDEX IF EXISTS idx_screen_time_user_date;
CREATE INDEX idx_screen_time_user_device_date ON screen_time_logs (user_id, device_id, recorded_at DESC);