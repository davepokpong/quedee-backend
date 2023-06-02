package domains

import (
	"context"

	"github.com/davepokpong/eticket-backend/dto"
	"github.com/davepokpong/eticket-backend/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserUseCase interface {
	SignIn(ctx context.Context, username string, password string) (models.User, string, error)
	Register(ctx context.Context, urd dto.UserRegisterDto, requesterRole string) error
	DeleteUser(ctx context.Context, username string) error
	EditUser(ctx context.Context, uud dto.UserUpdateDto) (models.User, error, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByUserName(ctx context.Context, username string) (models.User, error)
	ChangePassword(ctx context.Context, upud dto.UserPasswordUpdateDto) (models.User, error)
	GetAllByUserType(ctx context.Context, userType string) ([]models.User, error)
	GetUsersFromDateRangeAndRole(ctx context.Context, min string, max string, role string) ([]models.User, error)
	GetUsersPerWeek(ctx context.Context, min string, max string, role string) (dto.WeekDayStat, error)
	GetUsersPerYear(ctx context.Context, date string, role string) ([]models.User, error) 
	GetUsersPerMonthInAYear(ctx context.Context, max string, role string) (dto.StatPerMonth, error) 
	GetCustomerAndFilterByMembers(ctx context.Context) ([]dto.CountMemberOfUser, error)
	GetUserActivityListByUsername(ctx context.Context, username string) (dto.UserActivityListDto, error)
	FindActiveUserForLogin(ctx context.Context, username string) (models.User, error)
}

type UserRepository interface {
	FindUser(ctx context.Context, username string) (models.User, error)
	CreateUser(ctx context.Context, um models.User) error
	DeleteUser(ctx context.Context, username string) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	EditUser(ctx context.Context, username string, update interface{}) (models.User, error, error)
	GetAllByUserType(ctx context.Context, userType string) ([]models.User, error)
	GetAllUsersWithFilterSortedByMembers(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	GetUsersWithFilter(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	GetUsersPerYearFromRole(ctx context.Context, year int, role string) ([]models.User, error)
	CountUser(ctx context.Context, filter interface{}) (int64, error)
}
