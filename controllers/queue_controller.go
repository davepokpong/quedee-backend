package controllers

import (
	"net/http"

	"github.com/davepokpong/eticket-backend/domains"
	"github.com/davepokpong/eticket-backend/dto"
	"github.com/gin-gonic/gin"
)

type QueueController struct {
	queueUseCase domains.QueueUseCase
}

func NewQueueController(qu domains.QueueUseCase) *QueueController {
	return &QueueController{
		queueUseCase: qu,
	}
}

func (qc *QueueController) CreateQueue(c *gin.Context) {
	var body dto.CreateQueueDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	queue, err := qc.queueUseCase.CreateQueue(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"queue": queue,
		"status": true,
	})
}

func (qc *QueueController) DeleteQueue(c *gin.Context) {
	var body dto.QueueDeleteDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err := qc.queueUseCase.DeleteQueue(c, body.ID)
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

func (qc *QueueController) EditQueue(c *gin.Context) {
	var body dto.QueueUpdateDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	queue, err := qc.queueUseCase.EditQueue(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message1": err.Error(),
			// "message2": errr.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"queue": queue,
	})
}

func (qc *QueueController) GetAllQueues(c *gin.Context) {
	queues, err := qc.queueUseCase.GetAllQueue(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"queue": queues,
	})
}

func (qc *QueueController) GetQueueById(c *gin.Context){
	id := c.Param("id")
	queue, err := qc.queueUseCase.GetQueueById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"queue": queue,
	})
}

func (qc *QueueController) StartActivity(c *gin.Context) {
	var body dto.QueueDisableDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	queue, err := qc.queueUseCase.StartActivity(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"queue": queue,
	})
}

func (qc *QueueController) FindActiveQueueByUsername(c *gin.Context){
	username := c.Param("username")
	queues, err := qc.queueUseCase.FindActiveQueueByUsername(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"queue": queues,
	})
}

func (qc *QueueController) CancelQueue(c *gin.Context) {
	var body dto.QueueCancelDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	_, err := qc.queueUseCase.CancelQueue(c, body)
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

func (qc *QueueController) GetFirstQueue(c *gin.Context){
	username := c.Param("username")
	queue, err := qc.queueUseCase.CheckFirstQueue(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"queue": queue,
	})
}

func (qc *QueueController) GetActivityPerWeek(c *gin.Context){
	code := c.Param("code")
	max := c.Query("max")
	min := c.Query("min")
	stat, err := qc.queueUseCase.GetActivityPerWeek(c, min, max, code)
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

func (qc *QueueController) GetQueueFromDateRange(c *gin.Context){
	code := c.Param("code")
	max := c.Query("max")
	min := c.Query("min")
	stat, err := qc.queueUseCase.GetQueuesFromDateRange(c, min, max, code)
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

func (qc *QueueController) GetActivityFromYear(c *gin.Context){
	code := c.Param("code")
	max := c.Query("max")
	stat, err := qc.queueUseCase.GetActivityPerMonthInAYear(c, max, code)
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

func (qc *QueueController) GetCustomerNumberFilterByCode(c *gin.Context){
	code := c.Param("code")
	stat, err := qc.queueUseCase.GetCustomerNumberFilterByCode(c, code)
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

func (qc *QueueController) GetRatioOfQueueByCode(c *gin.Context){
	code := c.Param("code")
	stat, err := qc.queueUseCase.GetRatioOfQueueByCode(c, code)
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

func (qc *QueueController) NumberActivityPlayedTime(c *gin.Context){
	stat, err := qc.queueUseCase.NumberActivityPlayedTime(c)
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

func (qc *QueueController) ResetData(c *gin.Context){
	err := qc.queueUseCase.ClearData(c)
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

func (qc *QueueController) GetDataCountStat(c *gin.Context){
	stat, err := qc.queueUseCase.GetDataCountStat(c)
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

func (qc *QueueController) GetActiveQueueFilterByActivityCode(c *gin.Context){
	code := c.Param("code")
	queues, err := qc.queueUseCase.GetActiveQueueFilterByActivityCode(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"queues": queues,
	})
}

func (qc *QueueController) CreateQueueSpecificRound(c *gin.Context) {
	var body dto.CreateQueueSpecificRoundDto
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	queue, err := qc.queueUseCase.CreateQueueSpecificRound(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"queue": queue,
		"status": true,
	})
}

func (qc *QueueController) TemporaryClosedActivity(c *gin.Context){
	code := c.Param("code")
	err := qc.queueUseCase.TemporaryClosedActivity(c, code)
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