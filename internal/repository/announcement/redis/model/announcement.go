package modelredis

type Announcement struct {
	ID               string `redis:"id"`
	Title            string `redis:"title"`
	Description      string `redis:"description"`
	Address          string `redis:"address"`
	Date             string `redis:"date"`
	Contacts         string `redis:"contacts"`
	ModerationStatus string `redis:"moderation_status"`
	SearchedStatus   bool   `redis:"searched_status"`
	UserID           string `redis:"owner_id"`
}
