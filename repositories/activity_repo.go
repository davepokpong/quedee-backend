package repositories

import (
	"context"

	"github.com/davepokpong/eticket-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type activityRepository struct {
	db *mongo.Database
}

func NewActivityRepository(db *mongo.Database) activityRepository {
	return activityRepository{
		db: db,
	}
}

func (a activityRepository) collection() *mongo.Collection {
	return a.db.Collection(models.Activity{}.CollectionName())
}

func (a activityRepository) FindActivity(ctx context.Context, code string) (models.Activity, error) {
	var activity models.Activity
	activityCollection := a.collection()

	err := activityCollection.FindOne(ctx, bson.M{"code": code}).Decode(&activity)
	return activity, err
}

func (a activityRepository) CreateActivity(ctx context.Context, am models.Activity) error {
	activityCollection := a.collection()
	_, err := activityCollection.InsertOne(ctx, &am)
	return err
}

func (a activityRepository) DeleteActivity(ctx context.Context, code string) error {
	activityCollection := a.collection()
	_, err := activityCollection.DeleteOne(ctx, bson.M{"code": code})
	return err
}

func (a activityRepository) GetAllActivities(ctx context.Context) ([]models.Activity, error) {
	activityCollection := a.collection()
	var activities []models.Activity

	res, err := activityCollection.Find(ctx, bson.M{})
	for res.Next(ctx) {
		var activity models.Activity
		res.Decode(&activity) // ma tam eak tee
		activities = append(activities, activity)
	}

	return activities, err
}

func (a activityRepository) EditActivity(ctx context.Context, code string, update interface{}) (models.Activity, error) {
	var activity models.Activity
	activityCollection := a.collection()
	err := activityCollection.FindOneAndUpdate(ctx, bson.M{"code": code}, bson.M{"$set": update}).Decode(&activity)
	newActivity, _ := a.FindActivity(ctx, code)

	return newActivity, err
}

func (a activityRepository) GetActivitiesWithFilter(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.M{"createdAt": -1})
	collection := a.collection()
	return collection.Find(ctx, filter, &queryOptions)
}

func (a activityRepository) CountActivity(ctx context.Context, filter interface{}) (int64, error){
	Collection := a.collection()
	count, err := Collection.CountDocuments(ctx, filter)
	return count, err
}
