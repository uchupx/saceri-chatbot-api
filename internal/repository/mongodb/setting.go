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

func (r *SettingRepoMongodb) Create(ctx context.Context, data models.SettingModel) (*models.SettingModel, error) {
	collection := r.db.Database(databaseName).Collection(r.collectionName)

	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now

	if data.Id.IsZero() {
		data.Id = bson.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *SettingRepoMongodb) Update(ctx context.Context, data models.SettingModel) (*models.SettingModel, error) {
	collection := r.db.Database(databaseName).Collection(r.collectionName)

	// Set update timestamp
	data.UpdatedAt = time.Now()

	filter := bson.M{
		"_id": data.Id,
	}

	update := bson.M{
		"$set": bson.M{
			"value":      data.Value,
			"updated_at": data.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedData models.SettingModel
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &updatedData, nil
}

func (r *SettingRepoMongodb) GetByKey(ctx context.Context, key models.SettingKey) (*models.SettingModel, error) {
	collection := r.db.Database(databaseName).Collection(r.collectionName)

	filter := bson.M{
		"key": key,
	}

	var data models.SettingModel
	err := collection.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r *SettingRepoMongodb) GetAllSettings(ctx context.Context) ([]models.SettingModel, error) {
	collection := r.db.Database(databaseName).Collection(r.collectionName)

	filter := bson.M{}

	var settings []models.SettingModel
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
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
