CREATE TABLE IF NOT EXISTS preset_messages(
    id BIGSERIAL PRIMARY KEY,
    message TEXT NOT NULL UNIQUE DEFAULT 'message',
    price INT NOT NULL DEFAULT 1000,
    rarity VARCHAR(20) NOT NULL DEFAULT 'common',
    is_active BOOLEAN DEFAULT true,       
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);