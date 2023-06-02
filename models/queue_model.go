package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Queue struct {
	ID              primitive.ObjectID `bson:"_id" json:"_id"`
	Username        string             `bson:"username" json:"username"`
	ActivityCode    string             `bson:"activityCode" json:"activityCode"`
	ActivityName    []string           `bson:"activityName" json:"activityName"`
	ActivityPicture string             `bson:"activityPicture" json:"activityPicture"`
	QueueNo         int                `bson:"queueNo" json:"queueNo"`
	Round           int                `bson:"round" json:"round"`
	Disable         bool               `bson:"disable" json:"disable"`
	Status          string             `bson:"status" json:"status"`
	Size            int                `bson:"size" json:"size"`
	Star            int                `bson:"star" json:"star"`
	DiffRound       int                `bson:"diffRound" json:"diffRound"`
	Duration        int                `bson:"duration" json:"duration"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (Queue) CollectionName() string {
	return "queue"
}

func (e *Queue) MarshalBSON() ([]byte, error) {
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

	type ue Queue
	return bson.Marshal((*ue)(e))
}
