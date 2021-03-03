package models

import (
	"time"

	u "github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Model struct {
	UUID      u.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
