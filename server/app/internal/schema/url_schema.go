package schema

import (
	"time"

	"github.com/google/uuid"
)

type UrlRequest struct {
	Label string `json:"label" validate:"required"`
	Url   string `json:"url" validate:"url,required"`
	// Interval int    `json:"interval" validate:"required"`
	Active *bool `json:"active" validate:"required"`
}

type UrlResponse struct {
	ID          uuid.UUID  `json:"id"`
	Label       string     `json:"label"`
	URL         string     `json:"url"`
	Interval    int        `json:"interval"`
	Active      bool       `json:"active"`
	LastChecked *time.Time `json:"last_checked"`
	CreatedAt   time.Time  `json:"created_at"`
}
