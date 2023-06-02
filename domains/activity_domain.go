package domains

import (
	"context"

	"github.com/davepokpong/eticket-backend/dto"
	"github.com/davepokpong/eticket-backend/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActivityUseCase interface {
	CreateActivity(ctx context.Context, cad dto.CreateActivityDto) error
	DeleteActivity(ctx context.Context, code string) error
	GetAllActivities(ctx context.Context) ([]models.Activity, error)
	EditActivity(ctx context.Context, aud dto.ActivityUpdateDto) (models.Activity, error)
	GetActivityByCode(ctx context.Context, code string) (models.Activity, error)
	UpdateActivityComment(ctx context.Context, uac dto.UpdateActivityComment) (models.Activity, error, error)
	RatingCalculate(rating float32, commentNo int, inputRating float32) float32
	// UpdateActivityQueue(ctx context.Context, uaq dto.UpdateActivityQueue) (models.Activity, error)
	CheckWaitRound(ctx context.Context) ([]dto.ActivityWaitRoundDto, error)
}

type ActivityRepository interface {
	FindActivity(ctx context.Context, code string) (models.Activity, error)
	CreateActivity(ctx context.Context, am models.Activity) error
	DeleteActivity(ctx context.Context, code string) error
	GetAllActivities(ctx context.Context) ([]models.Activity, error)
	EditActivity(ctx context.Context, code string, update interface{}) (models.Activity, error)
	GetActivitiesWithFilter(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	CountActivity(ctx context.Context, filter interface{}) (int64, error)
}
