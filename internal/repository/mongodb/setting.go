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

type SettingRepoMongodb struct {
	db             *mongo.Client
	collectionName string
}

func (r *SettingRepoMongodb) Create(data models.SettingModel) (*models.SettingModel, error) {
	collection := r.db.Database(databaseName).Collection(r.collectionName)

	// Set timestamps
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now

	// Generate new ObjectID if not provided
	if data.Id.IsZero() {
		data.Id = bson.NewObjectID()
	}

	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *SettingRepoMongodb) Update(data models.SettingModel) (*models.SettingModel, error) {
	collection := r.db.Database(databaseName).Collection(r.collectionName)

	// Set update timestamp
	data.UpdatedAt = time.Now()

	filter := bson.M{
		"_id":       data.Id,
		"is_active": true,
	}

	update := bson.M{
		"$set": bson.M{
			"value":      data.Value,
			"updated_at": data.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedData models.SettingModel
	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &updatedData, nil
}

func (r *SettingRepoMongodb) GetByKey(key models.SettingKey) (*models.SettingModel, error) {
	collection := r.db.Database(databaseName).Collection(r.collectionName)

	filter := bson.M{
		"key":       key,
		"is_active": true,
	}

	var data models.SettingModel
	err := collection.FindOne(context.Background(), filter).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r *SettingRepoMongodb) GetAllSettings() ([]models.SettingModel, error) {
	collection := r.db.Database(databaseName).Collection(r.collectionName)

	filter := bson.M{"is_active": true}

	var settings []models.SettingModel
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var setting models.SettingModel
		err := cursor.Decode(&setting)
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return settings, nil
}

func NewSettingRepoMongodb(client *mongo.Client) *SettingRepoMongodb {
	return &SettingRepoMongodb{
		db:             client,
		collectionName: "settings",
	}
}
