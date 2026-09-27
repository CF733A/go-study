CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE INDEX IF NOT EXISTS idx_users_is_deleted ON users(is_deleted);

COMMENT ON INDEX idx_users_email IS 'Индекс для поиска по email';

COMMENT ON INDEX idx_users_is_deleted IS 'Индекс для фильтрации удалённых пользователей';