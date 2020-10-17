package models

import (
	"context"
	"fmt"

	"github.com/samderlust/spa_manager/utils/logger"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getIgnoreCaseSearch(search string) primitive.M {
	return primitive.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}
}

func saveEntity(entity interface{}, collection *mongo.Collection) (*primitive.ObjectID, *resterrors.RestError) {
	ctx := context.Background()
	res, err := collection.InsertOne(ctx, entity)
	if err != nil {
		return nil, resterrors.NewInternalServerError(fmt.Sprintf("Error saving entity, %s", err.Error()))
	}

	id := res.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func getMultipleEntities(filter map[string]interface{}, collection *mongo.Collection) (*mongo.Cursor, *resterrors.RestError) {
	ctx := context.Background()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, resterrors.NewBadRequestError("error finding entities")
	}
	return cursor, nil
}

func getEntity(entity interface{}, filter map[string]interface{}, collection *mongo.Collection) *resterrors.RestError {
	ctx := context.Background()
	res := collection.FindOne(ctx, filter)

	if err := res.Decode(entity); err != nil {
		return resterrors.NewInternalServerError("Error getting entity")
	}

	return nil
}

func updateEntity(entity interface{}, filter, update interface{}, collection *mongo.Collection) *resterrors.RestError {
	ctx := context.Background()
	res := collection.FindOneAndUpdate(ctx, filter, update)

	if err := res.Decode(&entity); err != nil {
		logger.Error("updating err", err)
		return resterrors.NewInternalServerError("error updating entity")
	}
	return nil
}

func deleteEntity(entity interface{}, filter map[string]interface{}, collection *mongo.Collection) *resterrors.RestError {
	ctx := context.Background()
	res := collection.FindOneAndDelete(ctx, filter)
	if err := res.Decode(entity); err != nil {
		logger.Error("err delete", err)
		return resterrors.NewInternalServerError("error deleting an entity")
	}
	return nil

}
