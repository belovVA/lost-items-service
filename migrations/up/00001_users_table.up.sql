CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                     name VARCHAR(100) NOT NULL,
                                     surname VARCHAR(100) NOT NULL,
                                     email VARCHAR(255) NOT NULL UNIQUE,
                                     phone VARCHAR(20) NOT NULL,
                                     password VARCHAR(255) NOT NULL,
                                     role VARCHAR(50) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
