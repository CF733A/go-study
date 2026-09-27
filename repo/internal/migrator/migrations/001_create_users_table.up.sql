CREATE TABLE IF NOT EXISTS users(

    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    name VARCHAR(255) NOT NULL,
     
    email VARCHAR(255) UNIQUE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    deleted_at TIMESTAMP NULL,
    
    is_deleted BOOLEAN DEFAULT FALSE
);

COMMENT ON TABLE users IS 'Таблица для хранения данных пользователей';