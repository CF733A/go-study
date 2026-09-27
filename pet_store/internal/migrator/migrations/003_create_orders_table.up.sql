CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    quantity INTEGER DEFAULT 1,
    ship_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'placed',
    complete BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_orders_pet FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE,
    CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_orders_pet_id ON orders(pet_id);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_complete ON orders(complete);

COMMENT ON TABLE orders IS 'Таблица заказов PetStore';
COMMENT ON INDEX idx_orders_pet_id IS 'Индекс для поиска заказов по питомцу';
COMMENT ON INDEX idx_orders_user_id IS 'Индекс для поиска заказов по пользователю';
COMMENT ON INDEX idx_orders_status IS 'Индекс для поиска по статусу заказа';
COMMENT ON INDEX idx_orders_complete IS 'Индекс для фильтрации завершенных заказов';