package main

import (
	"context"
	"fmt"

	"github.com/davepokpong/eticket-backend/controllers"
	"github.com/davepokpong/eticket-backend/repositories"
	"github.com/davepokpong/eticket-backend/usecases"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoUri string `mapstructure:"MONGO_URI"`
}

func main() {
	router := gin.Default()
	// configs.ConnectDB()
	ctx, _ := context.WithCancel(context.Background())

	var conf Config
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoUri))
	// mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:123456@localhost:27017"))
	if err != nil {
		panic(err)
	}

	db := mongoClient.Database("ETicket")

	userRepository := repositories.NewUserRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepository)

	activityRepository := repositories.NewActivityRepository(db)
	activityUseCase := usecases.NewActivityUseCase(activityRepository, userRepository)

	queueRepository := repositories.NewQueueRepository(db)
	queueUseCase := usecases.NewQueueUseCase(queueRepository, activityRepository, userRepository)

	//routes
	controllers.SetUpRoutes(router, userUseCase, activityUseCase, queueUseCase)

	router.Run()

	// s := gocron.NewScheduler(time.UTC)
	// s.Every(1).Day().At("10:30").Do(func(){
	// 	queueUseCase.ClearData(ctx)
	// })
}
