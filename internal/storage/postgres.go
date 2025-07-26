package storage

import (
	"context"
	"errors"
	"my-pastebin/internal/paste"
	"time"

	"gorm.io/gorm"
)

type DBMeter interface {
	IncDBQuery(operation string)
}

type Storage struct {
	db      *gorm.DB
	metrics DBMeter
}

func New(db *gorm.DB, metrics DBMeter) *Storage {
	return &Storage{db: db, metrics: metrics}
}

func (s *Storage) Save(ctx context.Context, p *paste.Paste) error {
	s.metrics.IncDBQuery("save")
	result := s.db.WithContext(ctx).Create(p)
	return result.Error
}

func (s *Storage) GetByID(ctx context.Context, id string) (*paste.Paste, error) {
	s.metrics.IncDBQuery(("get_by_id"))
	var p paste.Paste
	result := s.db.WithContext(ctx).First(&p, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &p, nil
}

func (s *Storage) DeleteExpired(ctx context.Context) (int64, error) {
	s.metrics.IncDBQuery("delete_expired")
	result := s.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&paste.Paste{})

	return result.RowsAffected, result.Error
}
