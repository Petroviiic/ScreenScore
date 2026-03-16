ALTER TABLE screen_time_logs 
ADD COLUMN device_id VARCHAR(255) NOT NULL DEFAULT 'default_device';