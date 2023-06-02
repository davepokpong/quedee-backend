package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/davepokpong/eticket-backend/domains"
	"github.com/davepokpong/eticket-backend/dto"
	"github.com/davepokpong/eticket-backend/models"
	"github.com/davepokpong/eticket-backend/utils"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo domains.UserRepository
}

func NewUserUseCase(ur domains.UserRepository) userUseCase {
	return userUseCase{
		userRepo: ur,
	}
}

func (uu userUseCase) SignIn(ctx context.Context, username string, password string) (models.User, string, error) {
	user, err := uu.FindActiveUserForLogin(ctx, username)
	if err != nil {
		return user, "", err
	}
	if username != user.Username {
		return user, "", fmt.Errorf("can't find user")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, "", fmt.Errorf("wrong password")
	}
	accessToken, err := createToken(user)
	return user, accessToken, err
}

func createToken(u models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = u.Username
	claims["role"] = u.Role
	// claims["user"] = u
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	accessToken, err := token.SignedString([]byte("secret")) //change later
	return accessToken, err
}

func (uu userUseCase) Register(ctx context.Context, urd dto.UserRegisterDto, requesterRole string) error {
	if !utils.DecideRole(requesterRole, urd.Role) {
		return errors.New("No Permission to register")
	}

	user, err := uu.userRepo.FindUser(ctx, urd.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if user.Username != "" && user.Disable == false {
		return errors.New("Number is already exist")
	}
	hashed, err := utils.GenerateHashPassword(urd.Password)
	if err != nil {
		return errors.New("Hash failed")
	}

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		FirstName: urd.FirstName,
		LastName:  urd.LastName,
		Username:  urd.Username,
		Email:     urd.Email,
		Phone:     urd.Phone,
		Password:  hashed,
		Role:      urd.Role,
		Members:   urd.Members,
		Star:      urd.Star,
		Disable:   false,
	}
	BASE_URL := "eticket-mailservice.kraikub.com"

	postURL := url.URL{
		Host:   BASE_URL,
		Path:   "/api/v1/verify-email",
		Scheme: "https",
	}

	// Prepare request body
	body := dto.PostUserDetail{
		To:       urd.Email,
		Code:     "https://quedeeproj.web.app/customer-signin",
		Name:     urd.FirstName,
		Password: urd.Password,
		Username: urd.Username,
		Lang:     "th",
	}
	fmt.Println(body)

	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(bodyBytes)
	// fmt.Println(reader)

	// Make HTTP POST request
	resp, err := http.Post(postURL.String(), "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}

	// Close response body
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Read response body
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
		log.Println("Error response. Status Code: ", resp.StatusCode)
	}

	return uu.userRepo.CreateUser(ctx, newUser)
}

func (uu userUseCase) DeleteUser(ctx context.Context, username string) error {
	err := uu.userRepo.DeleteUser(ctx, username)
	return err
}

func (uu userUseCase) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	users, err := uu.userRepo.GetAllUsers(ctx)
	return users, err
}

func (uu userUseCase) EditUser(ctx context.Context, uud dto.UserUpdateDto) (models.User, error, error) {
	update := bson.M{
		"firstname": uud.FirstName,
		"lastname":  uud.LastName,
		"username":  uud.Username,
		"email":     uud.Email,
		"phone":     uud.Phone,
		"members":   uud.Members,
		"star":      uud.Star,
		"disable":   uud.Disable,
	}
	user, err, errr := uu.userRepo.EditUser(ctx, uud.Username, update)
	return user, err, errr
}

func (uu userUseCase) GetUserByUserName(ctx context.Context, username string) (models.User, error) {
	user, err := uu.userRepo.FindUser(ctx, username)
	return user, err
}

func (uu userUseCase) ChangePassword(ctx context.Context, upud dto.UserPasswordUpdateDto) (models.User, error) {
	_, _, err := uu.SignIn(ctx, upud.Username, upud.Password)
	if err != nil {
		return models.User{}, err
	}
	newPassword, err := utils.GenerateHashPassword(upud.NewPassword)
	if err != nil {
		return models.User{}, err
	}

	update := bson.M{
		"password": newPassword,
	}
	user, err, errr := uu.userRepo.EditUser(ctx, upud.Username, update)
	if errr != nil {
		return models.User{}, errr
	}
	return user, err
}

func (uu userUseCase) GetAllByUserType(ctx context.Context, userType string) ([]models.User, error) {
	var users []models.User

	users, err := uu.userRepo.GetAllByUserType(ctx, userType)
	return users, err
}

func (uu userUseCase) GetUsersFromDateRangeAndRole(ctx context.Context, min string, max string, role string) ([]models.User, error) {
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
		"role": role,
	}

	cursor, err := uu.userRepo.GetUsersWithFilter(ctx, filter)
	if err != nil {
		return []models.User{}, err
	}
	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return []models.User{}, err
	}

	return users, err
}

func (uu userUseCase) GetUsersPerWeek(ctx context.Context, min string, max string, role string) (dto.WeekDayStat, error) {
	users, err := uu.GetUsersFromDateRangeAndRole(ctx, min, max, role)
	if err != nil {
		return dto.WeekDayStat{}, err
	}
	stat := dto.WeekDayStat{}

	for _, user := range users {
		// fmt.Println(user.CreatedAt.String())
		if user.CreatedAt.Weekday().String() == "Sunday"{
			stat.Sun += user.Members
		} else if user.CreatedAt.Weekday().String() == "Monday"{
			stat.Mon += user.Members
		} else if user.CreatedAt.Weekday().String() == "Tuesday"{
			stat.Tue += user.Members
		} else if user.CreatedAt.Weekday().String() == "Wednesday"{
			stat.Wed += user.Members
		} else if user.CreatedAt.Weekday().String() == "Thursday"{
			stat.Thu += user.Members
		} else if user.CreatedAt.Weekday().String() == "Friday"{
			stat.Fri += user.Members
		} else if user.CreatedAt.Weekday().String() == "Saturday"{
			stat.Sat += user.Members
		} 
	}

	return stat, err
}

func (uu userUseCase) GetUsersPerYear(ctx context.Context, date string, role string) ([]models.User, error) {
	dateISO, err := time.Parse(time.RFC3339, date)
    if err != nil {
        fmt.Println("Error parsing date:", err)
    }
	users, err := uu.userRepo.GetUsersPerYearFromRole(ctx, dateISO.Year(), role)
	if err != nil {
		return []models.User{}, err
	}

	return users, err
}

func (uu userUseCase) GetUsersPerMonthInAYear(ctx context.Context, max string, role string) (dto.StatPerMonth, error) {
	users, err := uu.GetUsersPerYear(ctx, max, role)
	if err != nil {
		return dto.StatPerMonth{}, err
	}
	stat := dto.StatPerMonth{}

	for _, user := range users {
		// fmt.Println(user.CreatedAt.String())
		if user.CreatedAt.Month().String() == "January"{
			stat.Jan += user.Members
		} else if user.CreatedAt.Month().String() == "February"{
			stat.Feb += user.Members
		} else if user.CreatedAt.Month().String() == "March"{
			stat.Mar += user.Members
		} else if user.CreatedAt.Month().String() == "April"{
			stat.Apr += user.Members
		} else if user.CreatedAt.Month().String() == "May"{
			stat.May += user.Members
		} else if user.CreatedAt.Month().String() == "June"{
			stat.Jun += user.Members
		} else if user.CreatedAt.Month().String() == "July"{
			stat.Jul += user.Members
		} else if user.CreatedAt.Month().String() == "August"{
			stat.Aug += user.Members
		} else if user.CreatedAt.Month().String() == "September"{
			stat.Sep += user.Members
		} else if user.CreatedAt.Month().String() == "October"{
			stat.Oct += user.Members
		} else if user.CreatedAt.Month().String() == "November"{
			stat.Nov += user.Members
		} else if user.CreatedAt.Month().String() == "December"{
			stat.Dec += user.Members
		} 
	}

	return stat, err
}

func (uu userUseCase) GetCustomerAndFilterByMembers(ctx context.Context) ([]dto.CountMemberOfUser, error) {

	filter := bson.M{
		"role": "customer",
	}

	cursor, err := uu.userRepo.GetAllUsersWithFilterSortedByMembers(ctx, filter)
	if err != nil {
		return []dto.CountMemberOfUser{}, err
	}
	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return []dto.CountMemberOfUser{}, err
	}
	tmp := 0
	var arr []dto.CountMemberOfUser
	max := users[0].Members 
	for index, user := range users {
		if user.Members == max{
			tmp += 1
		} else if user.Members < max{
			obj := dto.CountMemberOfUser{
				Member: max,
				Count: tmp,
			}
			arr = append(arr, obj)
			if len(arr) == 6 {
				break
			}
			max = user.Members
			tmp = 1
		}
		if len(users) == index+1 { //if the last one 
			objj := dto.CountMemberOfUser{
				Member: max,
				Count: tmp,
			}
			arr = append(arr, objj)
		}
	}

	return arr, err
}

func (uu userUseCase) GetUserActivityListByUsername(ctx context.Context, username string) (dto.UserActivityListDto, error){
	user, err := uu.userRepo.FindUser(ctx, username)
	if err != nil {
		return dto.UserActivityListDto{}, err
	}

	newDto := dto.UserActivityListDto{
		Activity: user.Activity,
	}
	return newDto, err
}

func (uu userUseCase) FindActiveUserForLogin(ctx context.Context, username string) (models.User, error){
	filter := bson.M{
		"disable": false,
		"username": username,
	}
	cursor, err := uu.userRepo.GetUsersWithFilter(ctx, filter)
	if err != nil {
		return models.User{}, err
	}
	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return models.User{}, err
	}
	if len(users) == 0{
		return models.User{}, err
	}
	return users[0], err
}
