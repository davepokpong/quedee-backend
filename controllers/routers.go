package controllers

import (
	"net/http"

	"github.com/davepokpong/eticket-backend/domains"
	"github.com/davepokpong/eticket-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine, uu domains.UserUseCase, au domains.ActivityUseCase, qu domains.QueueUseCase) {
	uc := UserController{userUseCase: uu}
	ac := ActivityController{activityUseCase: au}
	qc := QueueController{queueUseCase: qu}

	router.Use(middlewares.CORSMiddleware())
	authRoute(router, uc, qc)

	router.Use(middlewares.AuthorizeJWT())
	activityRoute(router, ac, qc)
	queueRoute(router, qc)
}

func authRoute(router *gin.Engine, uc UserController, qc QueueController) {

	router.POST("/auth/login", uc.SignIn)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Use(middlewares.AuthorizeJWT())
	router.POST("/auth/register", middlewares.AuthorizeAdminOrStaff(), uc.Register)
	router.POST("/auth/user/delete", middlewares.AuthorizeAdminOrStaff(), uc.DeleteUser)
	router.PUT("/auth/user", middlewares.AuthorizeAdminOrStaff(), uc.EditUser)
	router.GET("/auth/user/all", middlewares.AuthorizeAdminOrStaff(), uc.GetAllUsers)
	router.GET("/auth/user/:username", uc.GetUserByUserName)
	router.POST("/auth/user/changepassword", uc.ChangePassword)
	router.GET("/auth/user/all/:userType", middlewares.AuthorizeAdminOrStaff(), uc.GetAllByUserType)
	router.GET("/auth/user/perweek/customer", middlewares.AuthorizeAdminOrStaff(), uc.GetNumberOfCustomerFromWeek)
	router.GET("/auth/user/peryear/customer", middlewares.AuthorizeAdminOrStaff(), uc.GetNumberOfCustomerFromYear)
	router.GET("/auth/user/customer-filter-by-member", middlewares.AuthorizeAdminOrStaff(), uc.GetCustomerFilterByMembers)
	router.GET("/auth/user/activity/:username", uc.GetUserActivityListByUserName)
	router.GET("/auth/resetdata", middlewares.AuthorizeOnlyAdmin(), qc.ResetData)
	router.GET("/auth/datastat", middlewares.AuthorizeOnlyAdmin(), qc.GetDataCountStat)
}

func activityRoute(router *gin.Engine, ac ActivityController, qc QueueController) {

	router.GET("/activity/all", ac.GetAllActivities)
	router.GET("/activity/code/:code", ac.GetActivityByCode)
	router.POST("/activity/comment", ac.UpdateActivityComment)

	router.POST("/activity", middlewares.AuthorizeAdminOrStaff(), ac.CreateActivity)
	router.POST("/activity/delete", middlewares.AuthorizeAdminOrStaff(), ac.DeleteActivity)
	router.PUT("/activity", middlewares.AuthorizeAdminOrStaff(), ac.EditActivity)
	router.GET("/activity/waitround", ac.GetActivityWaitRound)
	router.GET("/activity/perweek/:code", middlewares.AuthorizeAdminOrStaff(), qc.GetActivityPerWeek)
	router.GET("/activity/peryear/:code", middlewares.AuthorizeAdminOrStaff(), qc.GetActivityFromYear)
	router.GET("/activity/queuesize/:code", middlewares.AuthorizeAdminOrStaff(), qc.GetCustomerNumberFilterByCode)
	router.GET("/activity/member-ratio/:code", middlewares.AuthorizeAdminOrStaff(), qc.GetRatioOfQueueByCode)
	router.GET("/activity/played-time/", middlewares.AuthorizeAdminOrStaff(), qc.NumberActivityPlayedTime)
	router.GET("/activity/temporary-close/:code", middlewares.AuthorizeAdminOrStaff(), qc.TemporaryClosedActivity)

}

func queueRoute(router *gin.Engine, qc QueueController) {

	router.POST("/queue", qc.CreateQueue)
	router.POST("/queue/cancel", qc.CancelQueue)
	router.POST("/queue/delete", qc.DeleteQueue)
	router.PUT("/queue", qc.EditQueue)
	router.GET("/queue/all", qc.GetAllQueues)
	router.GET("/queue/id/:id", qc.GetQueueById)
	router.POST("/queue/start", qc.StartActivity)
	router.GET("/queue/active/:username", qc.FindActiveQueueByUsername)
	router.GET("/queue/firstqueue/:username", qc.GetFirstQueue)
	router.GET("/queue/:code", qc.GetQueueFromDateRange)
	router.GET("/queue/round-per-activity/:code", middlewares.AuthorizeAdminOrStaff(), qc.GetActiveQueueFilterByActivityCode)
	router.POST("/queue/specific-round", qc.CreateQueueSpecificRound)

}
