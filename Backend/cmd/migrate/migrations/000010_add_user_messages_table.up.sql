CREATE TABLE IF NOT EXISTS user_messages(
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    message_id BIGINT REFERENCES preset_messages(id) ON DELETE CASCADE,
    purchased_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, message_id) 
);