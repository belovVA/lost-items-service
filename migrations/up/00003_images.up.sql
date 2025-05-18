CREATE TABLE IF NOT EXISTS announcement_images (
                                                   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                                   bytes BYTEA NOT NULL,
                                                   announcement_id UUID NOT NULL,
                                                   CONSTRAINT fk_announcement
                                                       FOREIGN KEY (announcement_id)
                                                           REFERENCES announcements(id)
                                                           ON DELETE CASCADE
);
