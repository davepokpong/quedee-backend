package controllers

import (
	"net/http"

	"github.com/davepokpong/eticket-backend/domains"
	"github.com/davepokpong/eticket-backend/dto"
	"github.com/gin-gonic/gin"
)

type ActivityController struct {
	activityUseCase domains.ActivityUseCase
}

func NewActivityController(au domains.ActivityUseCase) *ActivityController {
	return &ActivityController{
		activityUseCase: au,
	}
}

func (ac *ActivityController) CreateActivity(c *gin.Context) {
	var body dto.CreateActivityDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err := ac.activityUseCase.CreateActivity(c, body)
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

func (ac *ActivityController) DeleteActivity(c *gin.Context) {
	var body dto.ActivityDeleteDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err := ac.activityUseCase.DeleteActivity(c, body.Code)
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

func (ac *ActivityController) EditActivity(c *gin.Context) {
	var body dto.ActivityUpdateDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	activity, err := ac.activityUseCase.EditActivity(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message1": err.Error(),
			// "message2": errr.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"activity": activity,
	})
}

func (ac *ActivityController) GetAllActivities(c *gin.Context) {
	activities, err := ac.activityUseCase.GetAllActivities(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"activity": activities,
	})
}

func (ac *ActivityController) GetActivityByCode(c *gin.Context){
	code := c.Param("code")
	activity, err := ac.activityUseCase.GetActivityByCode(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"activity": activity,
	})
}

func (ac *ActivityController) UpdateActivityComment(c *gin.Context) {
	var body dto.UpdateActivityComment
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	activity, err, errr := ac.activityUseCase.UpdateActivityComment(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message1": err.Error(),
			"message2": errr.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"activity": activity,
	})
}

func (ac *ActivityController) GetActivityWaitRound(c *gin.Context){
	activities, err := ac.activityUseCase.CheckWaitRound(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": activities,
	})
}