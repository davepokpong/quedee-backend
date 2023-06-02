package usecases

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/davepokpong/eticket-backend/domains"
	"github.com/davepokpong/eticket-backend/dto"
	"github.com/davepokpong/eticket-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type queueUseCase struct {
	queueRepo    domains.QueueRepository
	activityRepo domains.ActivityRepository
	userRepo     domains.UserRepository
}

func NewQueueUseCase(qr domains.QueueRepository, ar domains.ActivityRepository, ur domains.UserRepository) queueUseCase {
	return queueUseCase{
		queueRepo:    qr,
		activityRepo: ar,
		userRepo:     ur,
	}
}

func (qu queueUseCase) CreateQueue(ctx context.Context, cqd dto.CreateQueueDto) (models.Queue, error) {
	activity, err := qu.activityRepo.FindActivity(ctx, cqd.ActivityCode)
	if err != nil {
		return models.Queue{}, errors.New("No activity in database")
	}
	user, err := qu.userRepo.FindUser(ctx, cqd.Username)
	if err != nil {
		return models.Queue{}, errors.New("No user in database")
	}

	if activity.Size < cqd.Size {
		return models.Queue{}, errors.New("Maximum Limit!")
	}
	diffStar := user.Star - cqd.Star
	fmt.Println(diffStar)
	if diffStar < 0 {
		return models.Queue{}, errors.New("Not enough star")
	}
	updateUser := bson.M{
		"star": diffStar,
	}

	activity.QueueNo += 1
	newID := primitive.NewObjectID()
	fmt.Println(newID.Hex())
	for index, listRound := range activity.AllRounds {
		if listRound.Space >= cqd.Size && listRound.Status == "wait" {
			fmt.Println("tam yu")
			listRound.QueueId = append(listRound.QueueId, newID.Hex())
			listRound.Space -= cqd.Size
			activity.AllRounds[index] = listRound
			activity.Round = index + 1

			updateActivity := bson.M{
				"queueNo":   activity.QueueNo,
				"round":     activity.Round,
				"allRounds": activity.AllRounds,
			}

			diffRound := activity.Round - activity.RoundNow
			newQueue := models.Queue{
				ID:              newID,
				Username:        cqd.Username,
				ActivityCode:    cqd.ActivityCode,
				ActivityName:    activity.Name,
				ActivityPicture: activity.Picture,
				QueueNo:         activity.QueueNo,
				Round:           activity.Round,
				Disable:         false,
				Status:          cqd.Status,
				Size:            cqd.Size,
				Star:            cqd.Star,
				DiffRound:       diffRound,
				Duration:        activity.Duration,
			}

			_, err = qu.activityRepo.EditActivity(ctx, activity.Code, updateActivity)
			if err != nil {
				return models.Queue{}, err
			}

			_, err, errr := qu.userRepo.EditUser(ctx, cqd.Username, updateUser)
			if err != nil {
				return models.Queue{}, err
			}
			if errr != nil {
				return models.Queue{}, errr
			}

			return newQueue, qu.queueRepo.CreateQueue(ctx, newQueue)
		}
	}

	numRound := len(activity.AllRounds)
	newListRound := models.ListRound{}
	newListRound.QueueId = append(newListRound.QueueId, newID.Hex())
	newListRound.Space = activity.Size - cqd.Size
	newListRound.Status = "wait"

	activity.Round = numRound + 1
	activity.AllRounds = append(activity.AllRounds, newListRound)

	updateActivity := bson.M{
		"queueNo":   activity.QueueNo,
		"round":     activity.Round,
		"allRounds": activity.AllRounds,
		"lastRound": activity.Round,
	}

	diffRound := activity.Round - activity.RoundNow
	newQueue := models.Queue{
		ID:              newID,
		Username:        cqd.Username,
		ActivityCode:    cqd.ActivityCode,
		ActivityName:    activity.Name,
		ActivityPicture: activity.Picture,
		QueueNo:         activity.QueueNo,
		Round:           activity.Round,
		Disable:         false,
		Status:          cqd.Status,
		Size:            cqd.Size,
		Star:            cqd.Star,
		DiffRound:       diffRound,
		Duration:        activity.Duration,
	}

	_, err = qu.activityRepo.EditActivity(ctx, activity.Code, updateActivity)
	if err != nil {
		return models.Queue{}, err
	}

	_, err, errr := qu.userRepo.EditUser(ctx, cqd.Username, updateUser)
	if err != nil {
		return models.Queue{}, err
	}
	if errr != nil {
		return models.Queue{}, errr
	}

	return newQueue, qu.queueRepo.CreateQueue(ctx, newQueue)

}

func (qu queueUseCase) DeleteQueue(ctx context.Context, id string) error {
	_, err := qu.queueRepo.FindQueue(ctx, id)
	if err != nil {
		return errors.New("No ducument in database")
	}
	errr := qu.queueRepo.DeleteQueue(ctx, id)
	return errr
}

func (qu queueUseCase) GetAllQueue(ctx context.Context) ([]models.Queue, error) {
	var queues []models.Queue
	queues, err := qu.queueRepo.GetAllQueue(ctx)

	return queues, err
}

func (qu queueUseCase) EditQueue(ctx context.Context, qud dto.QueueUpdateDto) (models.Queue, error) {
	update := bson.M{
		"username":     qud.Username,
		"activityCode": qud.ActivityCode,
		"queueNo":      qud.QueueNo,
		"round":        qud.Round,
		"disable":      qud.Disable,
		"status":       qud.Status,
		"size":         qud.Size,
	}
	queue, err := qu.queueRepo.EditQueue(ctx, qud.ID, update)
	return queue, err
}

func (qu queueUseCase) GetQueueById(ctx context.Context, id string) (models.Queue, error) {
	queue, err := qu.queueRepo.FindQueue(ctx, id)
	return queue, err
}

func (qu queueUseCase) DisableQueue(ctx context.Context, queueId string) (models.Queue, error) {
	update := bson.M{
		"disable": true,
		"status":  "done",
	}

	queue, err := qu.queueRepo.EditQueue(ctx, queueId, update)
	return queue, err
}

func (qu queueUseCase) Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (qu queueUseCase) StartActivity(ctx context.Context, qdd dto.QueueDisableDto) (models.Activity, error) {
	activity, err := qu.activityRepo.FindActivity(ctx, qdd.ActivityCode)
	if err != nil {
		return models.Activity{}, err
	}

	for _, queueId := range qdd.QueueId {
		_, err := qu.DisableQueue(ctx, queueId)
		if err != nil {
			return models.Activity{}, err
		}
		queue, err := qu.queueRepo.FindQueue(ctx, queueId)

		//check queue and remove it from other round
		thisRound := activity.RoundNow
		if queue.Round != thisRound {
			fmt.Println("other round")
			index := 0
			for i, qid := range activity.AllRounds[queue.Round-1].QueueId {
				if queueId == qid {
					index = i
					break
				}
			}
			activity.AllRounds[queue.Round-1].QueueId = remove(activity.AllRounds[queue.Round-1].QueueId, index)
			activity.AllRounds[queue.Round-1].Space += queue.Size

			updateActivity := bson.M{
				"allRounds": activity.AllRounds,
			}
			_, err := qu.activityRepo.EditActivity(ctx, qdd.ActivityCode, updateActivity)
			if err != nil {
				return models.Activity{}, err
			}
		}

		user, err := qu.userRepo.FindUser(ctx, queue.Username)
		updateUserActivity := models.ActivityList{
			Code:      queue.ActivityCode,
			Name:      activity.Name,
			Image:     activity.Picture,
			QueueSize: queue.Size,
			Star:      queue.Star,
			Status:    "done",
			QueueId:   queue.ID.Hex(),
		}
		user.Activity = append(user.Activity, updateUserActivity)
		updateUser := bson.M{
			"activity": user.Activity,
		}

		_, err, errr := qu.userRepo.EditUser(ctx, user.Username, updateUser)
		if errr != nil {
			return models.Activity{}, errr
		}
	}
	activity.AllRounds[activity.RoundNow-1].Status = "done"
	update := bson.M{
		"roundNow":  activity.RoundNow + 1,
		"allRounds": activity.AllRounds,
	}

	newActivity, err := qu.activityRepo.EditActivity(ctx, qdd.ActivityCode, update)
	for _, ar := range newActivity.AllRounds {
		if ar.Status == "wait" {
			for _, idq := range ar.QueueId {
				q, err := qu.queueRepo.FindQueue(ctx, idq)
				if err != nil {
					fmt.Println(err)
					return models.Activity{}, err
				}
				update := bson.M{
					"diffRound": q.DiffRound - 1,
				}
				_, err = qu.queueRepo.EditQueue(ctx, idq, update)
				if err != nil {
					fmt.Println(err)
					return models.Activity{}, err
				}
			}
		}
	}
	for _, qid := range activity.AllRounds[activity.RoundNow-1].QueueId {
		if qu.Contains(qdd.QueueId, qid) == false {
			update := bson.M{
				"disable": true,
				"status":  "notJoined",
			}
			_, err = qu.queueRepo.EditQueue(ctx, qid, update)
			if err != nil {
				return models.Activity{}, err
			}
			queue, err := qu.queueRepo.FindQueue(ctx, qid)
			if err != nil {
				return models.Activity{}, err
			}
			user, err := qu.userRepo.FindUser(ctx, queue.Username)
			if err != nil {
				return models.Activity{}, err
			}
			updateUserActivity := models.ActivityList{
				Code:      queue.ActivityCode,
				Name:      activity.Name,
				Image:     activity.Picture,
				QueueSize: queue.Size,
				Star:      queue.Star,
				Status:    "notJoined",
			}
			user.Activity = append(user.Activity, updateUserActivity)
			updateUser := bson.M{
				"activity": user.Activity,
			}
			_, err, errr := qu.userRepo.EditUser(ctx, user.Username, updateUser)
			// fmt.Println("yad laew")
			if err != nil {
				return models.Activity{}, err
			}
			if errr != nil {
				return models.Activity{}, errr
			}

		}
	}

	return newActivity, err
}

func (qu queueUseCase) FindActiveQueueByUsername(ctx context.Context, username string) ([]models.Queue, error) {
	var queues []models.Queue
	queues, err := qu.queueRepo.FindActiveQueueByUsername(ctx, username)

	return queues, err
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func (qu queueUseCase) CancelQueue(ctx context.Context, qcd dto.QueueCancelDto) (models.Queue, error) {
	queue, err := qu.queueRepo.FindQueue(ctx, qcd.QueueId)
	if err != nil {
		return models.Queue{}, err
	}
	user, err := qu.userRepo.FindUser(ctx, queue.Username)
	if err != nil {
		return models.Queue{}, err
	}
	activity, err := qu.activityRepo.FindActivity(ctx, queue.ActivityCode)
	if err != nil {
		return models.Queue{}, err
	}
	index := 0
	for i, queueId := range activity.AllRounds[queue.Round-1].QueueId {
		if qcd.QueueId == queueId {
			index = i
			activity.AllRounds[queue.Round-1].Space += queue.Size
			break
		}
	}
	activity.AllRounds[queue.Round-1].QueueId = remove(activity.AllRounds[queue.Round-1].QueueId, index)

	updateActivity := bson.M{
		"allRounds": activity.AllRounds,
	}

	newActivityList := models.ActivityList{
		Code:      activity.Code,
		Name:      activity.Name,
		Image:     activity.Picture,
		QueueSize: queue.Size,
		Star:      queue.Star,
		Status:    "cancel",
	}
	user.Activity = append(user.Activity, newActivityList)

	updateUser := bson.M{
		"star":     user.Star + queue.Star, //return star back to user
		"activity": user.Activity,
	}

	updateQueue := bson.M{
		"disable": true,
		"status":  "cancel",
	}

	newQueue, err := qu.queueRepo.EditQueue(ctx, qcd.QueueId, updateQueue)
	if err != nil {
		return models.Queue{}, err
	}
	_, err, errr := qu.userRepo.EditUser(ctx, queue.Username, updateUser)
	if errr != nil {
		return models.Queue{}, errr
	}
	_, err = qu.activityRepo.EditActivity(ctx, activity.Code, updateActivity)

	return newQueue, err
}

func (qu queueUseCase) CheckFirstQueue(ctx context.Context, username string) (models.Queue, error) {
	var queues []models.Queue
	queues, err := qu.queueRepo.FindActiveQueueByUsername(ctx, username)
	if err != nil {
		return models.Queue{}, err
	}
	if len(queues) == 0 {
		return models.Queue{}, err
	}

	if len(queues) == 1 {
		return queues[0], err
	}

	minTime := queues[0].DiffRound * queues[0].Duration
	minQueue := queues[0]
	for _, queue := range queues {
		if minTime > (queue.DiffRound * queue.Duration) {
			minTime = queue.DiffRound * queue.Duration
			minQueue = queue
		}
	}

	return minQueue, err
}

func (qu queueUseCase) FindDisableQueueAndStatusIsDone(ctx context.Context) ([]models.Queue, error) {
	filter := bson.M{
		"status":  "done",
		"disable": true,
	}

	cursor, err := qu.queueRepo.FindQueueWithFilter(ctx, filter)
	if err != nil {
		return []models.Queue{}, err
	}
	var queues []models.Queue
	err = cursor.All(ctx, &queues)
	if err != nil {
		return []models.Queue{}, err
	}
	return queues, err
}

func sortByCount(arr []dto.ActivityPlayedTime) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Count > arr[j].Count
	})
}

func (qu queueUseCase) NumberActivityPlayedTime(ctx context.Context) ([]dto.ActivityPlayedTime, error) {
	queues, err := qu.FindDisableQueueAndStatusIsDone(ctx)
	if err != nil {
		return []dto.ActivityPlayedTime{}, err
	}

	var arr []dto.ActivityPlayedTime
	for _, queue := range queues {
		if len(arr) == 0 {
			objChkLen := dto.ActivityPlayedTime{
				Code:  queue.ActivityCode,
				Name:  queue.ActivityName,
				Count: queue.Size,
			}
			arr = append(arr, objChkLen)
			continue
		}
		chk := false
		index := 0
		for i, d := range arr {
			if d.Code == queue.ActivityCode {
				index = i
				chk = true
				break
			}
		}
		if chk == true {
			arr[index].Count += queue.Size
			continue
		}
		obj := dto.ActivityPlayedTime{
			Code:  queue.ActivityCode,
			Name:  queue.ActivityName,
			Count: queue.Size,
		}
		arr = append(arr, obj)
	}

	sortByCount(arr)
	// fmt.Println(arr)
	var firstFive []dto.ActivityPlayedTime
	if len(arr) > 5 {
		firstFive = arr[:5]
	} else {
		firstFive = arr
	}

	return firstFive, err
}

func (qu queueUseCase) ClearActivityPerDay(ctx context.Context) error {
	activities, err := qu.activityRepo.GetAllActivities(ctx)
	if err != nil {
		fmt.Println("err find activity")
		return err
	}

	for _, activity := range activities {
		newListRound := models.ListRound{
			Space:  activity.Size,
			Status: "wait",
		}
		var newAllRounds []models.ListRound
		newAllRounds = append(newAllRounds, newListRound)
		update := bson.M{
			"queueNo":   0,
			"round":     1,
			"lastRound": 1,
			"roundNow":  1,
			"allRounds": newAllRounds,
		}
		_, err = qu.activityRepo.EditActivity(ctx, activity.Code, update)
		if err != nil {
			return err
		}
	}
	fmt.Println("clear activity")
	return err
}

func (qu queueUseCase) DisableCustomerPerDay(ctx context.Context) error {
	filter := bson.M{
		"disable": false,
		"role": "customer",
	}
	cursor, err := qu.userRepo.GetUsersWithFilter(ctx, filter)
	if err != nil {
		return err
	}
	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		fmt.Println("err find user")
		return err
	}
	if len(users) == 0{
		return err
	}

	for _, user := range users {
		if user.Role != "customer" {
			continue
		}
		update := bson.M{
			"username": user.Username + "-disable",
			"disable": true,
		}
		_, err, errr := qu.userRepo.EditUser(ctx, user.Username, update)
		if err != nil {
			return err
		}
		if errr != nil {
			return errr
		}
	}
	fmt.Println("clear user")
	return err
}

func (qu queueUseCase) DisableQueuePerDay(ctx context.Context) error {
	filter := bson.M{
		"disable": false,
	}
	cursor, err := qu.queueRepo.FindQueueWithFilter(ctx, filter)
	if err != nil {
		return err
	}
	var queues []models.Queue
	err = cursor.All(ctx, &queues)
	if err != nil {
		fmt.Println("err find queue")
		return err
	}
	if len(queues) == 0{
		return err
	}

	for _, queue := range queues {
		if queue.Disable == false {
			update := bson.M{
				"disable": true,
				"status":  "notJoined",
			}
			_, err = qu.queueRepo.EditQueue(ctx, queue.ID.Hex(), update)
			if err != nil {
				return err
			}
			activity, err := qu.activityRepo.FindActivity(ctx, queue.ActivityCode)
			if err != nil {
				return err
			}
			user, err := qu.userRepo.FindUser(ctx, queue.Username)
			if err != nil {
				return err
			}
			updateUserActivity := models.ActivityList{
				Code:      queue.ActivityCode,
				Name:      activity.Name,
				Image:     activity.Picture,
				QueueSize: queue.Size,
				Star:      queue.Star,
				Status:    "notJoined",
			}
			user.Activity = append(user.Activity, updateUserActivity)
			updateUser := bson.M{
				"activity": user.Activity,
			}
			_, err, errr := qu.userRepo.EditUser(ctx, user.Username, updateUser)
			if err != nil {
				return err
			}
			if errr != nil {
				return errr
			}
		}
	}
	return err
}

func (qu queueUseCase) ClearData(ctx context.Context) error {
	fmt.Println("start")
	err := qu.ClearActivityPerDay(ctx)
	if err != nil {
		return err
	}
	err = qu.DisableCustomerPerDay(ctx)
	if err != nil {
		return err
	}
	err = qu.DisableQueuePerDay(ctx)
	if err != nil {
		return err
	}
	return err
}

func (qu queueUseCase) GetQueuesFromDateRange(ctx context.Context, min string, max string, code string) ([]models.Queue, error) {
	minISO, err := time.Parse(time.RFC3339, min)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	// fmt.Println(minISO)
	maxISO, err := time.Parse(time.RFC3339, max)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	// fmt.Println(maxISO)
	filter := bson.M{
		"createdAt": bson.M{
			"$lte": maxISO,
			"$gte": minISO,
		},
		"activityCode": code,
		"status": "done",
	}

	cursor, err := qu.queueRepo.FindQueueWithFilter(ctx, filter)
	if err != nil {
		return []models.Queue{}, err
	}
	var queues []models.Queue
	err = cursor.All(ctx, &queues)
	if err != nil {
		return []models.Queue{}, err
	}

	return queues, err
}

func (qu queueUseCase) GetActivityPerWeek(ctx context.Context, min string, max string, code string) (dto.WeekDayStat, error) {
	queues, err := qu.GetQueuesFromDateRange(ctx, min, max, code)
	// fmt.Println(queues)
	if err != nil {
		return dto.WeekDayStat{}, err
	}
	stat := dto.WeekDayStat{}

	for _, queue := range queues {
		// fmt.Println(queue.CreatedAt.String())
		if queue.CreatedAt.Weekday().String() == "Sunday" {
			stat.Sun += queue.Size
		} else if queue.CreatedAt.Weekday().String() == "Monday" {
			stat.Mon += queue.Size
		} else if queue.CreatedAt.Weekday().String() == "Tuesday" {
			stat.Tue += queue.Size
		} else if queue.CreatedAt.Weekday().String() == "Wednesday" {
			stat.Wed += queue.Size
		} else if queue.CreatedAt.Weekday().String() == "Thursday" {
			stat.Thu += queue.Size
		} else if queue.CreatedAt.Weekday().String() == "Friday" {
			stat.Fri += queue.Size
		} else if queue.CreatedAt.Weekday().String() == "Saturday" {
			stat.Sat += queue.Size
		}
	}

	return stat, err
}

func (qu queueUseCase) GetQueuePerYear(ctx context.Context, date string, code string) ([]models.Queue, error) {
	dateISO, err := time.Parse(time.RFC3339, date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}
	queues, err := qu.queueRepo.GetQueuePerYearFromCode(ctx, dateISO.Year(), code)
	if err != nil {
		return []models.Queue{}, err
	}

	return queues, err
}

func (qu queueUseCase) GetActivityPerMonthInAYear(ctx context.Context, max string, code string) (dto.StatPerMonth, error) {
	queues, err := qu.GetQueuePerYear(ctx, max, code)
	if err != nil {
		return dto.StatPerMonth{}, err
	}
	stat := dto.StatPerMonth{}

	for _, queue := range queues {
		// fmt.Println(queue.CreatedAt.String())
		if queue.CreatedAt.Month().String() == "January" {
			stat.Jan += queue.Size
		} else if queue.CreatedAt.Month().String() == "February" {
			stat.Feb += queue.Size
		} else if queue.CreatedAt.Month().String() == "March" {
			stat.Mar += queue.Size
		} else if queue.CreatedAt.Month().String() == "April" {
			stat.Apr += queue.Size
		} else if queue.CreatedAt.Month().String() == "May" {
			stat.May += queue.Size
		} else if queue.CreatedAt.Month().String() == "June" {
			stat.Jun += queue.Size
		} else if queue.CreatedAt.Month().String() == "July" {
			stat.Jul += queue.Size
		} else if queue.CreatedAt.Month().String() == "August" {
			stat.Aug += queue.Size
		} else if queue.CreatedAt.Month().String() == "September" {
			stat.Sep += queue.Size
		} else if queue.CreatedAt.Month().String() == "October" {
			stat.Oct += queue.Size
		} else if queue.CreatedAt.Month().String() == "November" {
			stat.Nov += queue.Size
		} else if queue.CreatedAt.Month().String() == "December" {
			stat.Dec += queue.Size
		}
	}

	return stat, err
}

func (qu queueUseCase) GetCustomerNumberFilterByCode(ctx context.Context, code string) ([]dto.CountMemberOfUser, error) {
	filterQ := bson.M{
		"activityCode": code,
		"status": "done",
	}
	cursor, err := qu.queueRepo.FindQueueWithFilterSortedBySize(ctx, filterQ)
	if err != nil {
		return []dto.CountMemberOfUser{}, err
	}

	var queues []models.Queue
	err = cursor.All(ctx, &queues)
	if err != nil {
		return []dto.CountMemberOfUser{}, err
	}
	tmp := 0
	var arr []dto.CountMemberOfUser
	max := queues[0].Size
	for index, q := range queues {
		if q.Size == max {
			tmp += 1
		} else if q.Size < max {
			obj := dto.CountMemberOfUser{
				Member: max,
				Count:  tmp,
			}
			arr = append(arr, obj)
			if len(arr) == 6 {
				break
			}
			max = q.Size
			tmp = 1
		}
		if len(queues) == index+1 { //if the last one
			objj := dto.CountMemberOfUser{
				Member: max,
				Count:  tmp,
			}
			arr = append(arr, objj)
		}
	}

	return arr, err
}

func (qu queueUseCase) GetRatioOfQueueByCode(ctx context.Context, code string) (dto.CountQueueDto, error) {
	filterAll := bson.M{
		"disable": false,
	}
	all, err := qu.queueRepo.CountQueue(ctx, filterAll)
	if err != nil {
		return dto.CountQueueDto{}, err
	}
	filter := bson.M{
		"activityCode": code,
		"disable":      false,
	}
	acti, err := qu.queueRepo.CountQueue(ctx, filter)

	count := dto.CountQueueDto{
		All: all,
		One: acti,
	}

	return count, err
}

func (qu queueUseCase) GetDataCountStat(ctx context.Context) (dto.DataCountStatDto, error) {
	filterCustomer := bson.M{
		"role": "customer",
	}
	customer, err := qu.userRepo.CountUser(ctx, filterCustomer)
	if err != nil {
		return dto.DataCountStatDto{}, err
	}
	filterCustomerNow := bson.M{
		"role":    "customer",
		"disable": false,
	}
	customerNow, err := qu.userRepo.CountUser(ctx, filterCustomerNow)
	if err != nil {
		return dto.DataCountStatDto{}, err
	}
	filterStaff := bson.M{
		"role": bson.M{
			"$in": []string{"staff", "admin"},
		},
	}
	staff, err := qu.userRepo.CountUser(ctx, filterStaff)
	if err != nil {
		return dto.DataCountStatDto{}, err
	}
	filterStaffNow := bson.M{
		"role": bson.M{
			"$in": []string{"staff", "admin"},
		},
		"disable": false,
	}
	staffNow, err := qu.userRepo.CountUser(ctx, filterStaffNow)
	if err != nil {
		return dto.DataCountStatDto{}, err
	}
	filterActivity := bson.M{}
	activity, err := qu.activityRepo.CountActivity(ctx, filterActivity)
	if err != nil {
		return dto.DataCountStatDto{}, err
	}
	filterActivityNow := bson.M{
		"status": "open",
	}
	activityNow, err := qu.activityRepo.CountActivity(ctx, filterActivityNow)
	if err != nil {
		return dto.DataCountStatDto{}, err
	}

	obj := dto.DataCountStatDto{
		CustomerNow: customerNow,
		CustomerAll: customer,
		StaffNow:    staffNow,
		StaffAll:    staff,
		ActivityNow: activityNow,
		ActivityAll: activity,
	}

	return obj, err
}

func (qu queueUseCase) GetActiveQueueFilterByActivityCode(ctx context.Context, code string) ([]dto.ActiveQueueDto, error) {
	filter := bson.M{
		"activityCode": code,
		"disable":      false,
	}
	cursor, err := qu.queueRepo.FindQueueWithFilterSortedByRound(ctx, filter)
	if err != nil {
		return []dto.ActiveQueueDto{}, err
	}
	var queues []models.Queue
	err = cursor.All(ctx, &queues)
	if err != nil {
		return []dto.ActiveQueueDto{}, err
	}

	var arr []dto.ActiveQueueDto
	for _, queue := range queues {
		aqd := dto.ActiveQueueDto{
			QueueId:   queue.ID.Hex(),
			QueueSize: queue.Size,
			Round:     queue.Round,
		}
		arr = append(arr, aqd)
	}

	return arr, err
}

func (qu queueUseCase) CreateQueueSpecificRound(ctx context.Context, cqsr dto.CreateQueueSpecificRoundDto) (models.Queue, error) {
	activity, err := qu.activityRepo.FindActivity(ctx, cqsr.ActivityCode)
	if err != nil {
		return models.Queue{}, err
	}
	user, err := qu.userRepo.FindUser(ctx, cqsr.Username)
	if err != nil {
		return models.Queue{}, err
	}
	if cqsr.Size > activity.Size || cqsr.Size > activity.AllRounds[cqsr.Round-1].Space {
		return models.Queue{}, errors.New("Maximum limit!")
	}

	newID := primitive.NewObjectID()
	if activity.AllRounds[cqsr.Round-1].Status == "done" {
		return models.Queue{}, errors.New("This round is Done already")
	}
	activity.AllRounds[cqsr.Round-1].QueueId = append(activity.AllRounds[cqsr.Round-1].QueueId, newID.Hex())
	activity.AllRounds[cqsr.Round-1].Space -= cqsr.Size

	newQueue := models.Queue{
		ID:              newID,
		Username:        cqsr.Username,
		ActivityCode:    cqsr.ActivityCode,
		ActivityName:    activity.Name,
		ActivityPicture: activity.Picture,
		QueueNo:         activity.QueueNo + 1,
		Round:           cqsr.Round,
		Disable:         false,
		Status:          cqsr.Status,
		Size:            cqsr.Size,
		Star:            cqsr.Star,
		DiffRound:       cqsr.Round - activity.RoundNow,
		Duration:        activity.Duration,
	}

	//update activity
	updateActivity := bson.M{
		"allRounds": activity.AllRounds,
		"queueNo":   activity.QueueNo + 1,
	}
	_, err = qu.activityRepo.EditActivity(ctx, activity.Code, updateActivity)
	if err != nil {
		return models.Queue{}, err
	}
	//update user star
	diffStar := user.Star - cqsr.Star
	if diffStar < 0 {
		return models.Queue{}, errors.New("Not enough star")
	}
	updateUser := bson.M{
		"star": user.Star - cqsr.Star,
	}
	_, err, errr := qu.userRepo.EditUser(ctx, cqsr.Username, updateUser)
	if err != nil {
		return models.Queue{}, err
	}
	if errr != nil {
		return models.Queue{}, err
	}

	return newQueue, qu.queueRepo.CreateQueue(ctx, newQueue)
}

func (qu queueUseCase) TemporaryClosedActivity(ctx context.Context, code string) (error){
	filter := bson.M{
		"disable": false,
		"activityCode": code,
	}
	cursor, err := qu.queueRepo.FindQueueWithFilter(ctx, filter)
	if err != nil {
		return err
	}
	var queues []models.Queue
	err = cursor.All(ctx, &queues)
	if err != nil {
		return err
	}
	for _, queue := range queues {
		if queue.Disable == false {
			updateQ := bson.M{
				"status": "cancel",
				"disable": true,
			}
			_, err = qu.queueRepo.EditQueue(ctx, queue.ID.Hex(), updateQ)
			if err != nil {
				return err
			}
			newActivityList := models.ActivityList{
				Code: code,
				Name: queue.ActivityName,
				Image: queue.ActivityPicture,
				QueueId: queue.ID.Hex(),
				QueueSize: queue.Size,
				Star: queue.Star,
				Status: "cancel",
			}
			user, err := qu.userRepo.FindUser(ctx, queue.Username)
			if err != nil {
				return err
			}
			user.Activity = append(user.Activity, newActivityList)
			newStar := user.Star + queue.Star
			updateUser := bson.M{
				"star": newStar,
				"activity": user.Activity,
			}
			_, err, errr := qu.userRepo.EditUser(ctx, queue.Username, updateUser)
			if err != nil {
				return err
			}
			if errr != nil {
				return errr
			}
		}
	}
	activity, err := qu.activityRepo.FindActivity(ctx, code)
	if err != nil {
		return err
	}
	newListRound := models.ListRound{
		Space: activity.Size,
		Status: "wait",
	}
	newAllRound := []models.ListRound{}
	newAllRound = append(newAllRound, newListRound)
	updateActivity := bson.M {
		"allRounds": newAllRound,
	}
	_, err = qu.activityRepo.EditActivity(ctx, code, updateActivity)
	if err != nil {
		return err
	}
	return err

}