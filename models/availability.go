package models

import (
	"context"
	"time"

	"github.com/samderlust/spa_manager/resources"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Availability struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TechnicianID primitive.ObjectID `json:"technician,omimpty" bson:"technician,omitempty"`
	Date         time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	StartTime    time.Time          `json:"startTime,omitempty" bson:"startTime,omitempty"`
	EndTime      time.Time          `json:"endTime,omitempty" bson:"endTime,omitempty"`
}

func (a *Availability) Save() (*primitive.ObjectID, *error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := resources.Client.AvailabilityCollection().InsertOne(ctx, a)
	if err != nil {
		return nil, &err
	}
	id := res.InsertedID.(primitive.ObjectID)
	return &id, nil
}
func (a Availability) SaveMany(availabilities []Availability) (*[]primitive.ObjectID, *error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var docs []interface{}
	for _, v := range availabilities {
		docs = append(docs, v)
	}

	res, err := resources.Client.AvailabilityCollection().InsertMany(ctx, docs)
	if err != nil {
		return nil, &err
	}

	var ids []primitive.ObjectID
	for _, v := range res.InsertedIDs {
		ids = append(ids, v.(primitive.ObjectID))
	}

	return &ids, nil
}
