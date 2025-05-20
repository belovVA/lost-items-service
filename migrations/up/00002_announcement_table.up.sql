-- 1. Расширение для генерации uuid и полнотекстового поиска
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "unaccent";

-- 2. Таблица объявлений
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
                                             search_vector tsvector,  -- поле для поиска

                                             CONSTRAINT fk_user
                                                 FOREIGN KEY (user_id)
                                                     REFERENCES users(id)
                                                     ON DELETE CASCADE
);

-- 3. Индексы
CREATE INDEX IF NOT EXISTS idx_announcements_user_id ON announcements(user_id);

-- GIN-индекс по полнотекстовому полю
CREATE INDEX idx_announcements_search_vector
    ON announcements
        USING GIN (search_vector);

CREATE OR REPLACE FUNCTION announcements_search_vector_trigger() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
            to_tsvector('russian', unaccent(
                    coalesce(NEW.title, '') || ' ' ||
                    coalesce(NEW.description, '') || ' ' ||
                    coalesce(NEW.address, '') || ' ' ||
                    coalesce(NEW.contacts, '') || ' ' ||
                    coalesce(NEW.moderation_status, '') || ' ' ||
                    to_char(NEW.date, 'YYYY-MM-DD')
                                   ));
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- 5. Триггер на insert/update
CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
    ON announcements
    FOR EACH ROW EXECUTE FUNCTION announcements_search_vector_trigger();
