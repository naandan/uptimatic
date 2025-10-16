package services

import (
	"context"
	"errors"
	"net/http"
	"time"
	"uptimatic/internal/models"
	"uptimatic/internal/repositories"
	"uptimatic/internal/schema"
	"uptimatic/internal/utils"

	"gorm.io/gorm"
)

type URLService interface {
	Create(ctx context.Context, url *schema.UrlRequest, userID uint) (*schema.UrlResponse, *utils.AppError)
	Update(ctx context.Context, url *schema.UrlRequest, id uint) (*schema.UrlResponse, *utils.AppError)
	Delete(ctx context.Context, id uint) *utils.AppError
	FindByID(ctx context.Context, id uint) (*schema.UrlResponse, *utils.AppError)
	ListByUserID(ctx context.Context, userID uint, page, perPage int, active *bool, searchLabel string, sortBy string) ([]schema.UrlResponse, int, *utils.AppError)
	GetUptimeStats(ctx context.Context, urlID uint, mode string, offset int) ([]models.UptimeStat, *utils.AppError)
}

type urlService struct {
	db            *gorm.DB
	urlRepo       repositories.UrlRepository
	statusLogRepo repositories.StatusLogRepository
}

func NewUrlService(db *gorm.DB, urlRepo repositories.UrlRepository, statusLogRepo repositories.StatusLogRepository) URLService {
	return &urlService{db, urlRepo, statusLogRepo}
}

func (s *urlService) Create(ctx context.Context, url *schema.UrlRequest, userID uint) (*schema.UrlResponse, *utils.AppError) {
	// if !utils.ContainsInt(url.Interval) {
	// 	return nil, errors.New("invalid interval")
	// }

	urlModel := &models.URL{
		UserID:   userID,
		Label:    url.Label,
		URL:      url.Url,
		Interval: 300,
		Active:   *url.Active,
	}
	err := s.urlRepo.Create(ctx, s.db, urlModel)
	if err != nil {
		return nil, utils.InternalServerError("Error creating url", err)
	}
	return &schema.UrlResponse{
		ID:          urlModel.ID,
		Label:       urlModel.Label,
		URL:         urlModel.URL,
		Interval:    urlModel.Interval,
		Active:      urlModel.Active,
		LastChecked: urlModel.LastChecked,
		CreatedAt:   urlModel.CreatedAt,
	}, nil
}

func (s *urlService) Update(ctx context.Context, url *schema.UrlRequest, id uint) (*schema.UrlResponse, *utils.AppError) {
	urlModel, err := s.urlRepo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, utils.InternalServerError("Error finding url", err)
	}
	urlModel.Label = url.Label
	urlModel.URL = url.Url
	// urlModel.Interval = url.Interval
	urlModel.Active = *url.Active
	err = s.urlRepo.Update(ctx, s.db, urlModel)
	if err != nil {
		return nil, utils.InternalServerError("Error updating url", err)
	}
	return &schema.UrlResponse{
		ID:          urlModel.ID,
		Label:       urlModel.Label,
		URL:         urlModel.URL,
		Interval:    urlModel.Interval,
		Active:      urlModel.Active,
		LastChecked: urlModel.LastChecked,
		CreatedAt:   urlModel.CreatedAt,
	}, nil
}

func (s *urlService) Delete(ctx context.Context, id uint) *utils.AppError {
	urlModel, err := s.urlRepo.FindByID(ctx, s.db, id)
	if err != nil {
		return utils.NewAppError(http.StatusNotFound, utils.NotFound, "Url not found", err)
	}
	err = s.urlRepo.Delete(ctx, s.db, urlModel)
	if err != nil {
		return utils.InternalServerError("Error deleting url", err)
	}
	return nil
}

func (s *urlService) FindByID(ctx context.Context, id uint) (*schema.UrlResponse, *utils.AppError) {
	urlModel, err := s.urlRepo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, utils.NotFound, "Url not found", err)
	}
	return &schema.UrlResponse{
		ID:          urlModel.ID,
		Label:       urlModel.Label,
		URL:         urlModel.URL,
		Interval:    urlModel.Interval,
		Active:      urlModel.Active,
		LastChecked: urlModel.LastChecked,
		CreatedAt:   urlModel.CreatedAt,
	}, nil
}

func (s *urlService) ListByUserID(ctx context.Context, userID uint, page, perPage int, active *bool, searchLabel string, sortBy string) ([]schema.UrlResponse, int, *utils.AppError) {
	urls, count, err := s.urlRepo.ListByUserID(ctx, s.db, userID, page, perPage, active, searchLabel, sortBy)
	if err != nil {
		return nil, 0, utils.InternalServerError("Error listing urls", err)
	}
	var urlResponses []schema.UrlResponse
	for _, url := range urls {
		urlResponses = append(urlResponses, schema.UrlResponse{
			ID:          url.ID,
			Label:       url.Label,
			URL:         url.URL,
			Interval:    url.Interval,
			Active:      url.Active,
			LastChecked: url.LastChecked,
			CreatedAt:   url.CreatedAt,
		})
	}
	if len(urlResponses) == 0 {
		return []schema.UrlResponse{}, 0, nil
	}
	return urlResponses, count, nil
}

func (s *urlService) GetUptimeStats(ctx context.Context, urlID uint, mode string, offset int) ([]models.UptimeStat, *utils.AppError) {
	var targetDate time.Time

	switch mode {
	case "day":
		targetDate = time.Now().AddDate(0, 0, -offset)
	case "month":
		targetDate = time.Now().AddDate(0, -offset, 0)
	default:
		return nil, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid mode", nil)
	}

	logs, err := s.statusLogRepo.GetUptimeStats(ctx, s.db, urlID, mode, targetDate.UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.UptimeStat{}, nil
		}
		return nil, utils.InternalServerError("Error getting uptime stats", err)
	}

	return logs, nil
}
