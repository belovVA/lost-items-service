package modelpg

import "github.com/google/uuid"

type Image struct {
	ID    uuid.UUID `db:"id"`
	Bytes []byte    `db:"bytes"`
	AnnID uuid.UUID `db:"announcement_id"`
}
