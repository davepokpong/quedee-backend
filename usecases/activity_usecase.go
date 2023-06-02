package usecases

import (
	"context"
	"errors"

	"github.com/davepokpong/eticket-backend/domains"
	"github.com/davepokpong/eticket-backend/dto"
	"github.com/davepokpong/eticket-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type activityUseCase struct {
	activityRepo domains.ActivityRepository
	userRepo     domains.UserRepository
}

func NewActivityUseCase(ar domains.ActivityRepository, ur domains.UserRepository) activityUseCase {
	return activityUseCase{
		activityRepo: ar,
		userRepo:     ur,
	}
}

func (au activityUseCase) CreateActivity(ctx context.Context, cad dto.CreateActivityDto) error {
	activity, err := au.activityRepo.FindActivity(ctx, cad.Code)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if activity.Code != "" {
		return errors.New("Activity code is already exist")
	}
	newRoundList := models.ListRound{
		QueueId: []string{},
		Space:   cad.Size,
		Status:  "wait",
	}

	newActivity := models.Activity{
		ID:        primitive.NewObjectID(),
		Code:      cad.Code,
		Status:    cad.Status,
		Name:      cad.Name,
		Size:      cad.Size,
		Duration:  cad.Duration,
		Star:      cad.Star,
		Rating:    cad.Rating,
		Comment:   cad.Comment,
		Round:     1,
		RoundNow:  1,
		Position:  cad.Position,
		Picture:   cad.Picture,
		AllRounds: []models.ListRound{},
	}
	newActivity.AllRounds = append(newActivity.AllRounds, newRoundList)

	return au.activityRepo.CreateActivity(ctx, newActivity)
}

func (au activityUseCase) DeleteActivity(ctx context.Context, code string) error {
	_, err := au.activityRepo.FindActivity(ctx, code)
	if err != nil {
		return errors.New("No ducument in database")
	}
	errr := au.activityRepo.DeleteActivity(ctx, code)
	return errr
}

func (au activityUseCase) GetAllActivities(ctx context.Context) ([]models.Activity, error) {
	var activities []models.Activity

	activities, err := au.activityRepo.GetAllActivities(ctx)
	return activities, err
}

func (au activityUseCase) EditActivity(ctx context.Context, aud dto.ActivityUpdateDto) (models.Activity, error) {
	update := bson.M{
		"code":     aud.Code,
		"status":   aud.Status,
		"name":     aud.Name,
		"size":     aud.Size,
		"duration": aud.Duration,
		"star":     aud.Star,
		"rating":   aud.Rating,
		"comment":  aud.Comment,
		"position": aud.Position,
		"picture":  aud.Picture,
	}
	activity, err := au.activityRepo.EditActivity(ctx, aud.Code, update)
	return activity, err
}

func (au activityUseCase) GetActivityByCode(ctx context.Context, code string) (models.Activity, error) {
	activity, err := au.activityRepo.FindActivity(ctx, code)
	return activity, err
}

func (au activityUseCase) RatingCalculate(rating float32, commentNo int, inputRating float32) float32 {
	avgRating := ((rating * float32(commentNo)) + inputRating) / (float32(commentNo) + 1)

	return avgRating
}

func (au activityUseCase) UpdateActivityComment(ctx context.Context, uac dto.UpdateActivityComment) (models.Activity, error, error) {
	activity, err := au.activityRepo.FindActivity(ctx, uac.Code)
	if err != nil {
		return models.Activity{}, err, err
	}
	activity.Comment = append(activity.Comment, uac.Comment)

	avgRating := au.RatingCalculate(activity.Rating, activity.CommentNumber, uac.Comment.Rating)
	newCommentNumber := activity.CommentNumber + 1

	update := bson.M{
		"comment":       activity.Comment,
		"rating":        avgRating,
		"commentNumber": newCommentNumber,
	}

	activity, errr := au.activityRepo.EditActivity(ctx, activity.Code, update)
	if errr != nil {
		return models.Activity{}, errr, errr
	}

	user, err := au.userRepo.FindUser(ctx, uac.Comment.UserName)
	if err != nil {
		return models.Activity{}, err, err
	}
	for i, listActivity := range user.Activity {
		if uac.Comment.QueueId == listActivity.QueueId {
			user.Activity[i].CommentStatus = true
			updateUser := bson.M{
				"activity": user.Activity,
			}
			_, err, errr = au.userRepo.EditUser(ctx, uac.Comment.UserName, updateUser)
			if err != nil {
				return models.Activity{}, errr, errr
			}
			if errr != nil {
				return models.Activity{}, errr, errr
			}
		}
	}
	return activity, err, errr
}

func (au activityUseCase) CheckWaitRound(ctx context.Context) ([]dto.ActivityWaitRoundDto, error) {
	activities, err := au.activityRepo.GetAllActivities(ctx)
	if err != nil {
		return []dto.ActivityWaitRoundDto{}, err
	}

	var arr []dto.ActivityWaitRoundDto
	var roundNow int
	var lastRound int

	chkFirst := false
	
	for _, activity := range activities {
		for index, round := range activity.AllRounds {
			if round.Status == "wait" && chkFirst == false {
				roundNow = index + 1
				lastRound = roundNow
				chkFirst = true
				continue
			}
			if round.Status == "wait" && round.Space != activity.Size {
				lastRound = index + 1
			}
		}
		newDto := dto.ActivityWaitRoundDto{
			Activity:  activity,
			WaitRound: lastRound - roundNow,
		}
		arr = append(arr, newDto)
		roundNow = 0
		lastRound = 0
		chkFirst = false
	}
	return arr, err
}
