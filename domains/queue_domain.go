package domains

import (
	"context"

	"github.com/davepokpong/eticket-backend/dto"
	"github.com/davepokpong/eticket-backend/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type QueueUseCase interface {
	CreateQueue(ctx context.Context, cqd dto.CreateQueueDto) (models.Queue, error)
	DeleteQueue(ctx context.Context, id string) error
	GetAllQueue(ctx context.Context) ([]models.Queue, error)
	EditQueue(ctx context.Context, qud dto.QueueUpdateDto) (models.Queue, error)
	GetQueueById(ctx context.Context, id string) (models.Queue, error)
	DisableQueue(ctx context.Context, queueId string) (models.Queue, error)
	StartActivity(ctx context.Context, qdd dto.QueueDisableDto) (models.Activity, error)
	FindActiveQueueByUsername(ctx context.Context, username string) ([]models.Queue, error)
	CancelQueue(ctx context.Context, qcd dto.QueueCancelDto) (models.Queue, error)
	CheckFirstQueue(ctx context.Context, username string) (models.Queue, error)
	GetActivityPerWeek(ctx context.Context, min string, max string, code string) (dto.WeekDayStat, error)
	GetQueuesFromDateRange(ctx context.Context, min string, max string, code string) ([]models.Queue, error)
	GetActivityPerMonthInAYear(ctx context.Context, max string, code string) (dto.StatPerMonth, error)
	GetCustomerNumberFilterByCode(ctx context.Context, code string) ([]dto.CountMemberOfUser, error) 
	GetRatioOfQueueByCode(ctx context.Context, code string) (dto.CountQueueDto, error)
	NumberActivityPlayedTime(ctx context.Context) ([]dto.ActivityPlayedTime, error)
	ClearData(ctx context.Context) error
	GetDataCountStat(ctx context.Context) (dto.DataCountStatDto, error)
	GetActiveQueueFilterByActivityCode(ctx context.Context, code string) ([]dto.ActiveQueueDto, error)
	CreateQueueSpecificRound(ctx context.Context, cqsr dto.CreateQueueSpecificRoundDto) (models.Queue, error)
	TemporaryClosedActivity(ctx context.Context, code string) (error)
}

type QueueRepository interface {
	CreateQueue(ctx context.Context, qm models.Queue) error
	DeleteQueue(ctx context.Context, id string) error
	GetAllQueue(ctx context.Context) ([]models.Queue, error)
	EditQueue(ctx context.Context, id string, update interface{}) (models.Queue, error)
	FindQueue(ctx context.Context, id string) (models.Queue, error)
	FindActiveQueueByUsername(ctx context.Context, username string) ([]models.Queue, error)
	FindQueueWithFilter(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	GetQueuePerYearFromCode(ctx context.Context, year int, code string) ([]models.Queue, error)
	FindQueueWithFilterSortedBySize(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	CountQueue(ctx context.Context, filter interface{}) (int64, error)
	FindQueueWithFilterSortedByRound(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
}
