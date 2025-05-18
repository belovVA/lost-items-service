CREATE TABLE IF NOT EXISTS announcements (
                                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                             title VARCHAR(255) NOT NULL,
                                             description TEXT NOT NULL,
                                             address VARCHAR(255) NOT NULL,
                                             date TIMESTAMP NOT NULL,
                                             contacts VARCHAR(255) NOT NULL,
                                             searched_status BOOLEAN NOT NULL,
                                             moderation_status VARCHAR(50) NOT NULL,
                                             user_id UUID NOT NULL,
                                             CONSTRAINT fk_user
                                                 FOREIGN KEY (user_id)
                                                     REFERENCES users(id)
                                                     ON DELETE CASCADE
);

-- Индекс по user_id для ускорения поиска по владельцу
CREATE INDEX IF NOT EXISTS idx_announcements_user_id ON announcements(user_id);
