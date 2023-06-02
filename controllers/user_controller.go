package controllers

import (
	"net/http"

	"github.com/davepokpong/eticket-backend/domains"
	"github.com/davepokpong/eticket-backend/dto"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase domains.UserUseCase
}

func NewUserController(uu domains.UserUseCase) *UserController {
	return &UserController{
		userUseCase: uu,
	}
}

func (uc *UserController) SignIn(c *gin.Context) {
	body := new(dto.UserSignInDto)
	if err := c.BindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user, token, err := uc.userUseCase.SignIn(c, body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"accesstoken": token,
		"username": user.Username,
		"firstname": user.FirstName,
		"lastname": user.LastName,
		"email": user.Email,
		"members": user.Members,
		"star": user.Star,
		"role": user.Role,
		"phone": user.Phone,
		"_id": user.ID,
	})
}

func (uc *UserController) Register(c *gin.Context) {
	var body dto.UserRegisterDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err := uc.userUseCase.Register(c, body, c.GetString("role"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	var body dto.UserDeleteDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err := uc.userUseCase.DeleteUser(c, body.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func (uc *UserController) EditUser(c *gin.Context) {
	var body dto.UserUpdateDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user, err, errr := uc.userUseCase.EditUser(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message1": err.Error(),
			"message2": errr.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"user": user,
	})
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userUseCase.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"user": users,
	})
}

func (uc *UserController) GetUserByUserName(c *gin.Context){
	username := c.Param("username")
	user, err := uc.userUseCase.GetUserByUserName(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"user": user,
	})
}

func (uc *UserController) ChangePassword(c *gin.Context) {
	var body dto.UserPasswordUpdateDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user, err := uc.userUseCase.ChangePassword(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"user": user,
	})
}

func (uc *UserController) GetAllByUserType(c *gin.Context){
	userType := c.Param("userType")
	user, err := uc.userUseCase.GetAllByUserType(c, userType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"user": user,
	})
}

func (uc *UserController) GetNumberOfCustomerFromWeek(c *gin.Context){
	min := c.Query("min")
	max := c.Query("max")
	stat, err := uc.userUseCase.GetUsersPerWeek(c, min, max, "customer")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"stat": stat,
	})
}

func (uc *UserController) GetNumberOfCustomerFromYear(c *gin.Context){
	max := c.Query("max")
	stat, err := uc.userUseCase.GetUsersPerMonthInAYear(c, max, "customer")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"stat": stat,
	})
}

func (uc *UserController) GetCustomerFilterByMembers(c *gin.Context){
	data, err := uc.userUseCase.GetCustomerAndFilterByMembers(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"stat": data,
	})
}

func (uc *UserController) GetUserActivityListByUserName(c *gin.Context){
	username := c.Param("username")
	activityList, err := uc.userUseCase.GetUserActivityListByUsername(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"activity": activityList,
	})
}