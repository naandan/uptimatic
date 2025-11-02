package url

import (
	"context"
	"errors"
	"net/http"
	"time"
	"uptimatic/internal/models"
	"uptimatic/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URLService interface {
	Create(ctx context.Context, url *UrlRequest, userID uint) (*UrlResponse, *utils.AppError)
	Update(ctx context.Context, url *UrlRequest, id uuid.UUID) (*UrlResponse, *utils.AppError)
	Delete(ctx context.Context, id uuid.UUID) *utils.AppError
	FindByID(ctx context.Context, id uuid.UUID) (*UrlResponse, *utils.AppError)
	ListByUserID(ctx context.Context, userID uint, page, perPage int, active *bool, searchLabel string, sortBy string) ([]UrlResponse, int, *utils.AppError)
	GetUptimeStats(ctx context.Context, urlID uuid.UUID, mode, dateStr string) ([]models.UptimeStat, *utils.AppError)
}

type urlService struct {
	db            *gorm.DB
	urlRepo       UrlRepository
	statusLogRepo StatusLogRepository
}

func NewUrlService(db *gorm.DB, urlRepo UrlRepository, statusLogRepo StatusLogRepository) URLService {
	return &urlService{db, urlRepo, statusLogRepo}
}

func (s *urlService) Create(ctx context.Context, url *UrlRequest, userID uint) (*UrlResponse, *utils.AppError) {
	utils.Info(ctx, "Creating new URL", map[string]any{"user_id": userID, "label": url.Label, "url": url.Url})

	urlModel := &models.URL{
		UserID:   userID,
		PublicID: uuid.New(),
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
	return &UrlResponse{
		ID:          urlModel.PublicID,
		Label:       urlModel.Label,
		URL:         urlModel.URL,
		Interval:    urlModel.Interval,
		Active:      urlModel.Active,
		LastChecked: urlModel.LastChecked,
		CreatedAt:   urlModel.CreatedAt,
	}, nil
}

func (s *urlService) Update(ctx context.Context, url *UrlRequest, id uuid.UUID) (*UrlResponse, *utils.AppError) {
	utils.Info(ctx, "Updating URL", map[string]any{"url_id": id})

	urlModel, err := s.urlRepo.FindByPublicID(ctx, s.db, id)
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
	return &UrlResponse{
		ID:          urlModel.PublicID,
		Label:       urlModel.Label,
		URL:         urlModel.URL,
		Interval:    urlModel.Interval,
		Active:      urlModel.Active,
		LastChecked: urlModel.LastChecked,
		CreatedAt:   urlModel.CreatedAt,
	}, nil
}

func (s *urlService) Delete(ctx context.Context, id uuid.UUID) *utils.AppError {
	utils.Info(ctx, "Deleting URL", map[string]any{"url_id": id})

	urlModel, err := s.urlRepo.FindByPublicID(ctx, s.db, id)
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

func (s *urlService) FindByID(ctx context.Context, id uuid.UUID) (*UrlResponse, *utils.AppError) {
	utils.Info(ctx, "Fetching URL by ID", map[string]any{"url_id": id})

	urlModel, err := s.urlRepo.FindByPublicID(ctx, s.db, id)
	if err != nil {
		utils.Warn(ctx, "URL not found", map[string]any{"url_id": id})
		return nil, utils.NewAppError(http.StatusNotFound, utils.NotFound, "Url not found", err)
	}

	utils.Info(ctx, "URL fetched successfully", map[string]any{"url_id": id})
	return &UrlResponse{
		ID:          urlModel.PublicID,
		Label:       urlModel.Label,
		URL:         urlModel.URL,
		Interval:    urlModel.Interval,
		Active:      urlModel.Active,
		LastChecked: urlModel.LastChecked,
		CreatedAt:   urlModel.CreatedAt,
	}, nil
}

func (s *urlService) ListByUserID(ctx context.Context, userID uint, page, perPage int, active *bool, searchLabel string, sortBy string) ([]UrlResponse, int, *utils.AppError) {
	utils.Info(ctx, "Listing URLs by user", map[string]any{"user_id": userID, "page": page, "search": searchLabel})

	urls, count, err := s.urlRepo.ListByUserID(ctx, s.db, userID, page, perPage, active, searchLabel, sortBy)
	if err != nil {
		utils.Error(ctx, "Failed to list URLs", map[string]any{"user_id": userID, "err": err.Error()})
		return nil, 0, utils.InternalServerError("Error listing urls", err)
	}

	var responses []UrlResponse
	for _, url := range urls {
		responses = append(responses, UrlResponse{
			ID:          url.PublicID,
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
		return []UrlResponse{}, 0, nil
	}

	utils.Info(ctx, "URLs listed successfully", map[string]any{"user_id": userID, "count": len(responses)})
	return responses, count, nil
}

func (s *urlService) GetUptimeStats(ctx context.Context, urlID uuid.UUID, mode string, dateStr string) ([]models.UptimeStat, *utils.AppError) {
	url, err := s.urlRepo.FindByPublicID(ctx, s.db, urlID)
	if err != nil {
		utils.Warn(ctx, "URL not found", map[string]any{"url_id": urlID})
		return nil, utils.NewAppError(http.StatusNotFound, utils.NotFound, "Url not found", err)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		utils.Error(ctx, "Failed to load location", map[string]any{"err": err.Error()})
		return nil, utils.InternalServerError("Error loading location", err)
	}

	// Default: pakai tanggal hari ini kalau user tidak kasih parameter
	var targetDate time.Time
	if dateStr == "" {
		targetDate = time.Now().In(loc)
	} else {
		targetDate, err = time.ParseInLocation("2006-01-02", dateStr, loc)
		if err != nil {
			utils.Warn(ctx, "Invalid date format", map[string]any{"date": dateStr})
			return nil, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid date format (use YYYY-MM-DD)", err)
		}
	}

	var truncUnit string
	var startLocal, endLocal time.Time

	switch mode {
	case "day":
		truncUnit = "hour"
		startLocal = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, loc)
		endLocal = startLocal.Add(24 * time.Hour)

	case "month":
		truncUnit = "day"
		startLocal = time.Date(targetDate.Year(), targetDate.Month(), 1, 0, 0, 0, 0, loc)
		endLocal = startLocal.AddDate(0, 1, 0)

	case "year":
		truncUnit = "month"
		startLocal = time.Date(targetDate.Year(), 1, 1, 0, 0, 0, 0, loc)
		endLocal = startLocal.AddDate(1, 0, 0)

	default:
		utils.Warn(ctx, "Invalid mode for uptime stats", map[string]any{"mode": mode})
		return nil, utils.NewAppError(http.StatusBadRequest, utils.ValidationError, "Invalid mode", nil)
	}

	start := startLocal.UTC()
	end := endLocal.UTC()

	utils.Info(ctx, "Fetching uptime stats", map[string]any{
		"url_id":      urlID,
		"mode":        mode,
		"date":        targetDate.Format("2006-01-02"),
		"trunc_unit":  truncUnit,
		"start_local": startLocal,
		"end_local":   endLocal,
	})

	stats, err := s.statusLogRepo.GetUptimeStats(ctx, s.db, url.ID, truncUnit, start, end)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.UptimeStat{}, nil
		}
		utils.Error(ctx, "Failed to get uptime stats", map[string]any{"url_id": urlID, "err": err.Error()})
		return nil, utils.InternalServerError("Error getting uptime stats", err)
	}

	return stats, nil
}
