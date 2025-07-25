package storage

import (
	"context"
	"errors"
	"my-pastebin/internal/paste"
	"time"

	"gorm.io/gorm"
)

// Storage - наша структура для работы с БД
type Storage struct {
	db *gorm.DB
}

// New - конструктор для Storage
func New(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

// Save сохраняет новую пасту в БД
func (s *Storage) Save(ctx context.Context, p *paste.Paste) error {
	result := s.db.WithContext(ctx).Create(p)
	return result.Error
}

// GetByID получает пасту по ее уникальному ID
func (s *Storage) GetByID(ctx context.Context, id string) (*paste.Paste, error) {
	var p paste.Paste
	result := s.db.WithContext(ctx).First(&p, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound // Возвращаем специальную ошибку
		}
		return nil, result.Error
	}
	return &p, nil
}

// DeleteExpired удаляет все пасты с истекшим сроком хранения
func (s *Storage) DeleteExpired(ctx context.Context) (int64, error) {
	result := s.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).
		Delete(&paste.Paste{})

	// result.RowsAffected содержит кол-во удаленных строк
	return result.RowsAffected, result.Error
}
