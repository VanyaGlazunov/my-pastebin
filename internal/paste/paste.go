package paste

import (
	"time"
)

type Paste struct {
	ID        string    `gorm:"type:varchar(8);primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	Syntax    string    `gorm:"type:varchar(20);default:'text'"`
	CreatedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
