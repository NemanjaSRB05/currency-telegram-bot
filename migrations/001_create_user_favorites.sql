-- +goose Up
CREATE TABLE user_favorites (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Уникальный constraint чтобы один пользователь не мог добавить одну пару дважды
    UNIQUE(user_id, from_currency, to_currency)
);

CREATE INDEX idx_user_favorites_user_id ON user_favorites(user_id);
CREATE INDEX idx_user_favorites_currencies ON user_favorites(from_currency, to_currency);

-- +goose Down
DROP TABLE user_favorites;