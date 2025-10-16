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
	utils.Info(ctx, "Creating new URL", map[string]any{"user_id": userID, "label": url.Label, "url": url.Url})

	urlModel := &models.URL{
		UserID:   userID,
		Label:    url.Label,
		URL:      url.Url,
		Interval: 300,
		Active:   *url.Active,
	}

	err := s.urlRepo.Create(ctx, s.db, urlModel)
	if err != nil {
		utils.Error(ctx, "Failed to create URL", map[string]any{"user_id": userID, "err": err.Error()})
		return nil, utils.InternalServerError("Error creating url", err)
	}

	utils.Info(ctx, "URL created successfully", map[string]any{"url_id": urlModel.ID, "user_id": userID})
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
	utils.Info(ctx, "Updating URL", map[string]any{"url_id": id})

	urlModel, err := s.urlRepo.FindByID(ctx, s.db, id)
	if err != nil {
		utils.Error(ctx, "URL not found for update", map[string]any{"url_id": id, "err": err.Error()})
		return nil, utils.InternalServerError("Error finding url", err)
	}

	urlModel.Label = url.Label
	urlModel.URL = url.Url
	urlModel.Active = *url.Active

	err = s.urlRepo.Update(ctx, s.db, urlModel)
	if err != nil {
		utils.Error(ctx, "Failed to update URL", map[string]any{"url_id": id, "err": err.Error()})
		return nil, utils.InternalServerError("Error updating url", err)
	}

	utils.Info(ctx, "URL updated successfully", map[string]any{"url_id": id})
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
	utils.Info(ctx, "Deleting URL", map[string]any{"url_id": id})

	urlModel, err := s.urlRepo.FindByID(ctx, s.db, id)
	if err != nil {
		utils.Warn(ctx, "URL not found for deletion", map[string]any{"url_id": id})
		return utils.NewAppError(http.StatusNotFound, utils.NotFound, "Url not found", err)
	}

	err = s.urlRepo.Delete(ctx, s.db, urlModel)
	if err != nil {
		utils.Error(ctx, "Failed to delete URL", map[string]any{"url_id": id, "err": err.Error()})
		return utils.InternalServerError("Error deleting url", err)
	}

	utils.Info(ctx, "URL deleted successfully", map[string]any{"url_id": id})
	return nil
}

func (s *urlService) FindByID(ctx context.Context, id uint) (*schema.UrlResponse, *utils.AppError) {
	utils.Info(ctx, "Fetching URL by ID", map[string]any{"url_id": id})

	urlModel, err := s.urlRepo.FindByID(ctx, s.db, id)
	if err != nil {
		utils.Warn(ctx, "URL not found", map[string]any{"url_id": id})
		return nil, utils.NewAppError(http.StatusNotFound, utils.NotFound, "Url not found", err)
	}

	utils.Info(ctx, "URL fetched successfully", map[string]any{"url_id": id})
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
	utils.Info(ctx, "Listing URLs by user", map[string]any{"user_id": userID, "page": page, "search": searchLabel})

	urls, count, err := s.urlRepo.ListByUserID(ctx, s.db, userID, page, perPage, active, searchLabel, sortBy)
	if err != nil {
		utils.Error(ctx, "Failed to list URLs", map[string]any{"user_id": userID, "err": err.Error()})
		return nil, 0, utils.InternalServerError("Error listing urls", err)
	}

	var responses []schema.UrlResponse
	for _, url := range urls {
		responses = append(responses, schema.UrlResponse{
			ID:          url.ID,
			Label:       url.Label,
			URL:         url.URL,
			Interval:    url.Interval,
			Active:      url.Active,
			LastChecked: url.LastChecked,
			CreatedAt:   url.CreatedAt,
		})
	}

	if len(responses) == 0 {
		utils.Warn(ctx, "No URLs found for user", map[string]any{"user_id": userID})
		return []schema.UrlResponse{}, 0, nil
	}

	utils.Info(ctx, "URLs listed successfully", map[string]any{"user_id": userID, "count": len(responses)})
	return responses, count, nil
}

func (s *urlService) GetUptimeStats(ctx context.Context, urlID uint, mode string, offset int) ([]models.UptimeStat, *utils.AppError) {
	utils.Info(ctx, "Fetching uptime stats", map[string]any{"url_id": urlID, "mode": mode, "offset": offset})

	var targetDate time.Time
	switch mode {
	case "day":
		targetDate = time.Now().AddDate(0, 0, -offset)
	case "month":
		targetDate = time.Now().AddDate(0, -offset, 0)
	default:
		utils.Warn(ctx, "Invalid mode for uptime stats", map[string]any{"mode": mode})
		return nil, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid mode", nil)
	}

	stats, err := s.statusLogRepo.GetUptimeStats(ctx, s.db, urlID, mode, targetDate.UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Warn(ctx, "No uptime stats found", map[string]any{"url_id": urlID})
			return []models.UptimeStat{}, nil
		}
		utils.Error(ctx, "Failed to get uptime stats", map[string]any{"url_id": urlID, "err": err.Error()})
		return nil, utils.InternalServerError("Error getting uptime stats", err)
	}

	utils.Info(ctx, "Uptime stats fetched successfully", map[string]any{"url_id": urlID, "count": len(stats)})
	return stats, nil
}
