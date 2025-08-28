
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

type EventRepoMongodb struct {
    db             *mongo.Client
    collectionName string
}

func (r *EventRepoMongodb) Create(ctx context.Context, event models.EventModel) (*models.EventModel, error) {
    collection := r.db.Database(databaseName).Collection(r.collectionName)

    event.CreatedAt = time.Now()
    event.UpdatedAt = nil

    result, err := collection.InsertOne(ctx, event)
    if err != nil {
        return nil, err
    }

    event.Id = result.InsertedID.(bson.ObjectID)
    return &event, nil
}

func (r *EventRepoMongodb) Update(ctx context.Context, event models.EventModel) (*models.EventModel, error) {
    collection := r.db.Database(databaseName).Collection(r.collectionName)

    now := time.Now()
    event.UpdatedAt = &now

    filter := bson.M{"_id": event.Id}
    update := bson.M{"$set": event}

    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return nil, err
    }

    if result.MatchedCount == 0 {
        return nil, errors.New("event not found")
    }

    return &event, nil
}

func (r *EventRepoMongodb) Delete(ctx context.Context, id string) error {
    collection := r.db.Database(databaseName).Collection(r.collectionName)

    objectID, err := bson.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    now := time.Now()
    filter := bson.M{"_id": objectID}
    update := bson.M{"$set": bson.M{"deleted_at": now}}

    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return errors.New("event not found")
    }

    return nil
}

func (r *EventRepoMongodb) GetById(ctx context.Context, id string) (*models.EventModel, error) {
    collection := r.db.Database(databaseName).Collection(r.collectionName)

    objectID, err := bson.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    filter := bson.M{
        "_id":        objectID,
        "deleted_at": bson.M{"$exists": false},
    }

    var event models.EventModel
    err = collection.FindOne(ctx, filter).Decode(&event)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, errors.New("event not found")
        }
        return nil, err
    }

    return &event, nil
}

func (r *EventRepoMongodb) GetAllEvents(ctx context.Context, limit, offset int) ([]models.EventModel, error) {
    collection := r.db.Database(databaseName).Collection(r.collectionName)

    filter := bson.M{"deleted_at": bson.M{"$exists": false}}
    opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var events []models.EventModel
    if err = cursor.All(ctx, &events); err != nil {
        return nil, err
    }

    return events, nil
}

func NewEventRepoMongodb(client *mongo.Client) *EventRepoMongodb {
    return &EventRepoMongodb{
        db:             client,
        collectionName: "events",
    }
}
