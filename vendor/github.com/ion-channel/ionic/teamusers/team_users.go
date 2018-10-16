package teamusers

import (
	"time"
)

// TeamUser is a representation of an Ion Channel Team User relationship within the system
type TeamUser struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"team_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Status    string    `json:"status"`
	Role      string    `json:"role"`
}
