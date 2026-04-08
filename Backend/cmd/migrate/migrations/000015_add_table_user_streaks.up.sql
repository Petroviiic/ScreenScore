CREATE TABLE IF NOT EXISTS user_streaks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT NOT NULL DEFAULT 0, 
    all_time_high INT NOT NULL DEFAULT 0, 
    shield_count INT NOT NULL DEFAULT 0, 
    week_number INT NOT NULL, 
    year_number INT NOT NULL, 
    last_week_average NUMERIC(10, 2) DEFAULT 120.00, 
    last_updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);