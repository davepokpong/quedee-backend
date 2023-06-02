package repositories

import (
	"context"
	"log"

	"github.com/davepokpong/eticket-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type queueRepository struct {
	db *mongo.Database
}

func NewQueueRepository(db *mongo.Database) queueRepository {
	return queueRepository{
		db: db,
	}
}

func (q queueRepository) collection() *mongo.Collection {
	return q.db.Collection(models.Queue{}.CollectionName())
}

func (q queueRepository) FindQueue(ctx context.Context, id string) (models.Queue, error) {
	var queue models.Queue
	queueCollection := q.collection()
	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
		return models.Queue{}, err
	}

	err = queueCollection.FindOne(ctx, bson.M{"_id": newId}).Decode(&queue)
	return queue, err
}

func (q queueRepository) CreateQueue(ctx context.Context, qm models.Queue) error {
	queueCollection := q.collection()
	_, err := queueCollection.InsertOne(ctx, &qm)
	return err
}

func (q queueRepository) DeleteQueue(ctx context.Context, id string) error {
	queueCollection := q.collection()
	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
		return err
	}

	_, err = queueCollection.DeleteOne(ctx, bson.M{"_id": newId})
	return err
}

func (q queueRepository) GetAllQueue(ctx context.Context) ([]models.Queue, error) {
	queueCollection := q.collection()
	var queues []models.Queue

	res, err := queueCollection.Find(ctx, bson.M{})
	for res.Next(ctx) {
		var queue models.Queue
		res.Decode(&queue) // ma tam eak tee
		queues = append(queues, queue)
	}

	return queues, err
}

func (q queueRepository) EditQueue(ctx context.Context, id string, update interface{}) (models.Queue, error) {
	var queue models.Queue
	queueCollection := q.collection()
	newId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
		return models.Queue{}, err
	}
	err = queueCollection.FindOneAndUpdate(ctx, bson.M{"_id": newId}, bson.M{"$set": update}).Decode(&queue)
	newQueue, _ := q.FindQueue(ctx, id)

	return newQueue, err
}

func (q queueRepository) FindActiveQueueByUsername(ctx context.Context, username string) ([]models.Queue, error) {
	queueCollection := q.collection()
	var queues []models.Queue

	res, err := queueCollection.Find(ctx, bson.M{"username": username, "disable": false})
	for res.Next(ctx) {
		var queue models.Queue
		res.Decode(&queue) // ma tam eak tee
		queues = append(queues, queue)
	}

	return queues, err
}

func (q queueRepository) FindDisableQueueByUsername(ctx context.Context, username string) ([]models.Queue, error) {
	queueCollection := q.collection()
	var queues []models.Queue

	res, err := queueCollection.Find(ctx, bson.M{"username": username, "disable": true})
	for res.Next(ctx) {
		var queue models.Queue
		res.Decode(&queue)
		queues = append(queues, queue)
	}

	return queues, err
}

func (q queueRepository) FindQueueWithFilter(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.M{"createdAt": -1})
	collection := q.collection()
	return collection.Find(ctx, filter, &queryOptions)
}

func (q queueRepository) GetQueuePerYearFromCode(ctx context.Context, year int, code string) ([]models.Queue, error) {
	queueCollection := q.collection()
	var queues []models.Queue

	res, err := queueCollection.Find(ctx, bson.M{"activityCode": code, "status": "done"})
	for res.Next(ctx) {
		var queue models.Queue
		res.Decode(&queue)
		if queue.CreatedAt.Year() == year{
			queues = append(queues, queue)
		}
	}

	return queues, err
}

func (q queueRepository) FindQueueWithFilterSortedBySize(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.M{"size": -1})
	collection := q.collection()
	return collection.Find(ctx, filter, &queryOptions)
}

func (q queueRepository) CountQueue(ctx context.Context, filter interface{}) (int64, error){
	qCollection := q.collection()
	count, err := qCollection.CountDocuments(ctx, filter)
	return count, err
}

func (q queueRepository) FindQueueWithFilterSortedByRound(ctx context.Context, filter interface{}) (*mongo.Cursor, error){
	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.M{"round": -1})
	collection := q.collection()
	return collection.Find(ctx, filter, &queryOptions)
}

