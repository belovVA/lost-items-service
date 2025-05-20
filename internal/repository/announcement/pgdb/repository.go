package pgdb

import (
	"context"

	"github.com/google/uuid"
	"lost-items-service/internal/db/pgxdb"
	"lost-items-service/internal/model"
)

const (
	annTable                  = "announcements"
	annIDColumn               = "id"
	annTitleColumn            = "title"
	annDescColumn             = "description"
	annAddressColumn          = "address"
	annDateColumn             = "date"
	annContactsColumn         = "contacts"
	annSearchedStatusColumn   = "searched_status"
	annModerationStatusColumn = "moderation_status"
	userIDColumn              = "user_id"
)

type AnnPGRepository interface {
	AddAnn(ctx context.Context, ann *model.Announcement) (uuid.UUID, error)
	GetAnnByID(ctx context.Context, id uuid.UUID) (*model.Announcement, error)
	GetListAnnouncement(ctx context.Context, info *model.InfoSetting) ([]*model.Announcement, error)
}

type annRepo struct {
	DB pgxdb.DB
}

func NewRepository(db pgxdb.DB) AnnPGRepository {
	return &annRepo{
		DB: db,
	}
}
