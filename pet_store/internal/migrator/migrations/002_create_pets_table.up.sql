CREATE TABLE IF NOT EXISTS pets (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id BIGINT,
    category_name VARCHAR(100),
    photo_urls TEXT[] DEFAULT '{}',
    tags JSONB DEFAULT '[]',
    status VARCHAR(20) DEFAULT 'available',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pets_status ON pets(status);
CREATE INDEX IF NOT EXISTS idx_pets_category ON pets(category_id);
CREATE INDEX IF NOT EXISTS idx_pets_name ON pets(name);

COMMENT ON TABLE pets IS 'Таблица питомцев PetStore';
COMMENT ON INDEX idx_pets_status IS 'Индекс для поиска по статусу (available, pending, sold)';
COMMENT ON INDEX idx_pets_category IS 'Индекс для поиска по категории';
COMMENT ON INDEX idx_pets_name IS 'Индекс для поиска по имени питомца';