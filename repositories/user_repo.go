package repositories

import (
	"context"

	"github.com/davepokpong/eticket-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) userRepository {
	return userRepository{
		db: db,
	}
}

func (u userRepository) collection() *mongo.Collection {
	return u.db.Collection(models.User{}.CollectionName())
}

func (u userRepository) FindUser(ctx context.Context, username string) (models.User, error) {
	var user models.User
	userCollection := u.collection()

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return user, err
}

func (u userRepository) CreateUser(ctx context.Context, um models.User) error {
	userCollection := u.collection()
	_, err := userCollection.InsertOne(ctx, &um)
	return err
}

func (u userRepository) DeleteUser(ctx context.Context, username string) error {
	userCollection := u.collection()
	_, err := userCollection.DeleteOne(ctx, bson.M{"username": username})
	return err
}

func (u userRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	userCollection := u.collection()
	var users []models.User

	res, err := userCollection.Find(ctx, bson.M{})
	for res.Next(ctx) {
		var user models.User
		res.Decode(&user) // ma tam eak tee
		users = append(users, user)
	}

	return users, err
}

func (u userRepository) EditUser(ctx context.Context, username string, update interface{}) (models.User, error, error) {
	var user models.User
	userCollection := u.collection()
	err := userCollection.FindOneAndUpdate(ctx, bson.M{"username": username}, bson.M{"$set": update}).Decode(&user)
	newUser, errr := u.FindUser(ctx, username)

	return newUser, err, errr
}

func (u userRepository) GetAllByUserType(ctx context.Context, userType string) ([]models.User, error) {
	userCollection := u.collection()
	var users []models.User

	res, err := userCollection.Find(ctx, bson.M{"role": userType})
	for res.Next(ctx) {
		var user models.User
		res.Decode(&user) // ma tam eak tee
		users = append(users, user)
	}

	return users, err
}

func (u userRepository) GetAllUsersWithFilterSortedByMembers(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.M{"members": -1})
	collection := u.collection()
	return collection.Find(ctx, filter, &queryOptions)
}

func (u userRepository) GetUsersWithFilter(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.M{"createdAt": -1})
	collection := u.collection()
	return collection.Find(ctx, filter, &queryOptions)
}

func (u userRepository) GetUsersPerYearFromRole(ctx context.Context, year int, role string) ([]models.User, error) {
	userCollection := u.collection()
	var users []models.User

	res, err := userCollection.Find(ctx, bson.M{ "role": role })
	for res.Next(ctx) {
		var user models.User
		res.Decode(&user)
		if user.CreatedAt.Year() == year{
			users = append(users, user)
		}
	}

	return users, err
}

func (u userRepository) CountUser(ctx context.Context, filter interface{}) (int64, error){
	Collection := u.collection()
	count, err := Collection.CountDocuments(ctx, filter)
	return count, err
}