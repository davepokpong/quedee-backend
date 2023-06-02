package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserComment struct {
	UserName 	  string  			 `bson:"username" json:"username"`
	Rating        float32 			 `bson:"rating" json:"rating"`
	Text          string  			 `bson:"text" json:"text"`
	CreatedAt	  time.Time			 `bson:"createdAt" json:"createdAt"`
	QueueId 	  string			 `bson:"queueId" json:"queueId"`
}

type ListRound struct {
	QueueId 	  []string 			 `bson:"queueId" json:"queueId"`
	Space		  int				 `bson:"space" json:"space"`
	Status		  string			 `bson:"status" json:"status"`	
}

type Activity struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Code          string             `bson:"code" json:"code"`
	Status        string             `bson:"status" json:"status"`
	Name          []string           `bson:"name" json:"name"`
	Size          int                `bson:"size" json:"size"`
	Duration      int                `bson:"duration" json:"duration"`
	Star          int                `bson:"star" json:"star"`
	Rating        float32            `bson:"rating" json:"rating"`
	Comment       []UserComment      `bson:"comment" json:"comment"`
	CommentNumber int                `bson:"commentNumber" json:"commentNumber"`
	Position      []float32          `bson:"position" json:"position"`
	Picture       string             `bson:"picture" json:"picture"`
	QueueNo       int                `bson:"queueNo" json:"queueNo"`
	Round         int                `bson:"round" json:"round"`
	LastRound     int                `bson:"lastRound" json:"lastRound"`
	RoundNow      int                `bson:"roundNow" json:"roundNow"`
	AllRounds	  []ListRound		 `bson:"allRounds" json:"allRounds"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt	  time.Time		 `bson:"updatedAt" json:"updatedAt"`
}

func (Activity) CollectionName() string {
	return "activities"
}

func (e *Activity) MarshalBSON() ([]byte, error) {
	if e.ID.IsZero() {
		e.ID = primitive.NewObjectID()
	}

	thaiLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
    	fmt.Println(err)
	}
	now := time.Now().In(thaiLocation)
	if e.CreatedAt.IsZero() {
		e.CreatedAt = now
	}
	e.UpdatedAt = now

	type ue Activity
	return bson.Marshal((*ue)(e))
}
