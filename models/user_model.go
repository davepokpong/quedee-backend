package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityList struct {
	Code			string			   `bson:"code" json:"code"`
	Name			[]string		   `bson:"name" json:"name"`
	Image			string			   `bson:"image" json:"image"`
	QueueId			string			   `bson:"queueId" json:"queueId"`
	QueueSize  		int			   	   `bson:"queueSize" json:"queueSize"`
	Star			int				   `bson:"star" json:"star"`
	Status			string			   `bson:"status" json:"status"`
	CommentStatus	bool			   `bson:"commentStatus" json:"commentStatus"`	
}	

type User struct {
	ID        		primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName  	 	string             `bson:"firstname" json:"firstname"`
	LastName		string			   `bson:"lastname" json:"lastname"`
	Username		string 			   `bson:"username" json:"username"`
	Email			string			   `bson:"email" json:"email"`
	Phone			string			   `bson:"phone" json:"phone"`
	Password		string			   `bson:"password" json:"password"`
	Role			string			   `bson:"role" json:"role"`
	Members			int				   `bson:"members" json:"members"`
	Star 			int				   `bson:"star" json:"star"`
	Activity		[]ActivityList 	   `bson:"activity" json:"activity"`
	Disable			bool			   `bson:"disable" json:"disable"`
	CreatedAt 		time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt		time.Time		   `bson:"updatedAt" json:"updatedAt"`
}

func (User) CollectionName() string {
	return "users"
}

func (e *User) MarshalBSON() ([]byte, error) {
	
	if e.ID.IsZero() {
		e.ID = primitive.NewObjectID()
	}

	now := time.Now().UTC()
	if e.CreatedAt.IsZero() {
		e.CreatedAt = now
		fmt.Println(e.CreatedAt)
	}
	e.UpdatedAt = now

	type ue User
	return bson.Marshal((*ue)(e))
}
