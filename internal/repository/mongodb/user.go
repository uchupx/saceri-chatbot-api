package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/uchupx/saceri-chatbot-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepoMongodb struct {
	db *mongo.Client
}

const (
	databaseName   = "saceri_chatbot"
	collectionName = "users"
)

func (r *UserRepoMongodb) GetUser(ctx context.Context, id string) (*models.UserModel, error) {
	collection := r.db.Database(databaseName).Collection(collectionName)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	//
	filter := bson.M{
		"_id":       objectID,
		"is_active": true,
	}

	var user models.UserModel
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoMongodb) GetUserByOauthID(ctx context.Context, oauthID string) (*models.UserModel, error) {
	collection := r.db.Database(databaseName).Collection(collectionName)

	filter := bson.M{
		"oauth_id":  oauthID,
		"is_active": true,
	}

	var user models.UserModel
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoMongodb) GetUserByUsername(ctx context.Context, username string) (*models.UserModel, error) {
	collection := r.db.Database(databaseName).Collection(collectionName)

	filter := bson.M{
		"username":  username,
		"is_active": true,
	}

	var user models.UserModel
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoMongodb) CreateUser(ctx context.Context, user models.UserModel) (*models.UserModel, error) {
	collection := r.db.Database(databaseName).Collection(collectionName)

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsActive = true

	// Generate new ObjectID if not provided
	if user.Id.IsZero() {
		user.Id = bson.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoMongodb) UpdateUser(ctx context.Context, user models.UserModel) (*models.UserModel, error) {
	collection := r.db.Database(databaseName).Collection(collectionName)

	// Set update timestamp
	user.UpdatedAt = time.Now()

	filter := bson.M{
		"_id":       user.Id,
		"is_active": true,
	}

	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"updated_at": user.UpdatedAt,
		},
	}

	// Only update password if it's provided
	if user.Password != "" {
		update["$set"].(bson.M)["password"] = user.Password
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedUser models.UserModel
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &updatedUser, nil
}

func (r *UserRepoMongodb) DeleteUser(ctx context.Context, id string) error {
	collection := r.db.Database(databaseName).Collection(collectionName)

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":       objectID,
		"is_active": true,
	}

	// Soft delete - set is_active to false
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments // User not found
	}

	return nil
}

func (r *UserRepoMongodb) GetAllUsers(ctx context.Context, keyword *string, limit, offset int) ([]models.UserModel, error) {
	collection := r.db.Database(databaseName).Collection(collectionName)

	filter := bson.M{}
	if keyword != nil && *keyword != "" {
		filter["$or"] = []bson.M{
			{"username": bson.M{"$regex": *keyword, "$options": "i"}},
			{"name": bson.M{"$regex": *keyword, "$options": "i"}},
		}
	}

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.M{"created_at": -1}) // Sort by newest first

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.UserModel
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepoMongodb) CountUsers(ctx context.Context) (int64, error) {
	collection := r.db.Database(databaseName).Collection(collectionName)

	filter := bson.M{"is_active": true}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func NewUserRepoMongodb(db *mongo.Client) *UserRepoMongodb {
	return &UserRepoMongodb{
		db: db,
	}
}
