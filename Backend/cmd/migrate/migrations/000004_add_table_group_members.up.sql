CREATE TABLE IF NOT EXISTS group_members(
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (group_id, user_id)
);