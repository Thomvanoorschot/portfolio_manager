package entities

import (
	"github.com/google/uuid"
	"time"
)

type EntityBase struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
