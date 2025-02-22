package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinnaserwan/crm-be/services/feedback/internal/domain/model"
	"gorm.io/gorm"
)

type stationRepository struct {
	db *gorm.DB
}

func NewStationRepository(db *gorm.DB) *stationRepository {
	return &stationRepository{db: db}
}

func (r *stationRepository) List(ctx context.Context) ([]model.Station, error) {
	var stations []model.Station
	err := r.db.WithContext(ctx).Find(&stations).Error
	return stations, err
}

func (r *stationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Station, error) {
	var station model.Station
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&station).Error
	if err != nil {
		return nil, err
	}
	return &station, nil
}

func (r *stationRepository) Create(ctx context.Context, station *model.Station) error {
	return r.db.WithContext(ctx).Create(station).Error
}
