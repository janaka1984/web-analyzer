package repository

import (
	"context"
	"fmt"

	"github.com/janaka/web-analyzer/internal/config"
	"github.com/janaka/web-analyzer/internal/domain"
	"github.com/janaka/web-analyzer/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnalysisRepository interface {
	Save(ctx context.Context, a *domain.Analysis) error
	ListRecent(ctx context.Context, limit int) ([]domain.Analysis, error)
}

type analysisRepo struct{ db *gorm.DB }

func NewAnalysisRepository(cfg *config.Config, log logger.Logger) AnalysisRepository {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&domain.Analysis{}); err != nil {
		panic(err)
	}
	log.Info("DB connected & migrated (analysis table)")
	return &analysisRepo{db: db}
}

func (r *analysisRepo) Save(ctx context.Context, a *domain.Analysis) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *analysisRepo) ListRecent(ctx context.Context, limit int) ([]domain.Analysis, error) {
	var rows []domain.Analysis
	if err := r.db.WithContext(ctx).Order("id DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
