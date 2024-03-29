package repository

import (
	"context"

	"github.com/insabelter/IWS_GO/models"
)

type Repository interface {
	GetAllFeedbacks(ctx context.Context) ([]models.Feedback, error)
	GetFeedback(ctx context.Context, ID string) (models.Feedback, error)
	CreateFeedback(ctx context.Context, in models.Feedback) (models.Feedback, error)
	DeleteFeedback(ctx context.Context, ID string) error
}
