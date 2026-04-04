CREATE TABLE IF NOT EXISTS weekly_reward_logs (
    id BIGSERIAL PRIMARY KEY,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE, 
    week_number INT NOT NULL,                                       
    year_year INT NOT NULL,                                      
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT unique_group_weekly_reward UNIQUE (group_id, week_number, year_year)
);

CREATE INDEX idx_weekly_reward_lookup ON weekly_reward_logs (week_number, year_year);