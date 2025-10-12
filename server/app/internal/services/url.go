package services

import (
	"fmt"
	"time"
	"uptimatic/internal/models"
	"uptimatic/internal/repositories"
	"uptimatic/internal/schema"

	"gorm.io/gorm"
)

type URLService interface {
	Create(url *schema.UrlRequest, userID uint) (*schema.UrlResponse, error)
	Update(url *schema.UrlRequest, id uint) (*schema.UrlResponse, error)
	Delete(id uint) error
	FindByID(id uint) (*schema.UrlResponse, error)
	ListByUserID(userID uint, page, perPage int) ([]schema.UrlResponse, int, error)
	GetUptimeStats(urlID uint, mode string, offset int) ([]models.UptimeStat, error)
}

type urlService struct {
	db            *gorm.DB
	urlRepo       repositories.UrlRepository
	statusLogRepo repositories.StatusLogRepository
}

func NewUrlService(db *gorm.DB, urlRepo repositories.UrlRepository, statusLogRepo repositories.StatusLogRepository) URLService {
	return &urlService{db, urlRepo, statusLogRepo}
}

func (s *urlService) Create(url *schema.UrlRequest, userID uint) (*schema.UrlResponse, error) {
	// if !utils.ContainsInt(url.Interval) {
	// 	return nil, errors.New("invalid interval")
	// }

	urlModel := &models.URL{
		UserID:   userID,
		Label:    url.Label,
		URL:      url.Url,
		Interval: 300,
		Active:   url.Active,
	}
	err := s.urlRepo.Create(s.db, urlModel)
	if err != nil {
		return nil, err
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

func (s *urlService) Update(url *schema.UrlRequest, id uint) (*schema.UrlResponse, error) {
	urlModel, err := s.urlRepo.FindByID(s.db, id)
	if err != nil {
		return nil, err
	}
	urlModel.Label = url.Label
	urlModel.URL = url.Url
	// urlModel.Interval = url.Interval
	urlModel.Active = url.Active
	err = s.urlRepo.Update(s.db, urlModel)
	if err != nil {
		return nil, err
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

func (s *urlService) Delete(id uint) error {
	urlModel, err := s.urlRepo.FindByID(s.db, id)
	if err != nil {
		return err
	}
	return s.urlRepo.Delete(s.db, urlModel)
}

func (s *urlService) FindByID(id uint) (*schema.UrlResponse, error) {
	urlModel, err := s.urlRepo.FindByID(s.db, id)
	if err != nil {
		return nil, err
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

func (s *urlService) ListByUserID(userID uint, page, perPage int) ([]schema.UrlResponse, int, error) {
	urls, count, err := s.urlRepo.ListByUserID(s.db, userID, page, perPage)
	if err != nil {
		return nil, 0, err
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
	return urlResponses, count, nil
}

func (s *urlService) GetUptimeStats(urlID uint, mode string, offset int) ([]models.UptimeStat, error) {
	var targetDate time.Time

	switch mode {
	case "day":
		targetDate = time.Now().AddDate(0, 0, -offset)
	case "month":
		targetDate = time.Now().AddDate(0, -offset, 0)
	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}

	if _, err := s.statusLogRepo.GetByID(s.db, urlID); err != nil {
		return nil, err
	}

	return s.statusLogRepo.GetUptimeStats(s.db, urlID, mode, targetDate.UTC())
}
