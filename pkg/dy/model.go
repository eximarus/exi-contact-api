package dy

import "time"

type GuestbookEntry struct {
	Email     string
	Message   string
	Name      string
	Company   *string
	CreatedAt time.Time
}
