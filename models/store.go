package models

import (
	"context"

	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string               `json:"name,omitempty" bson:"name,omitempty"`
	Address      float64              `json:"address,omitempty" bson:"address,omitempty"`
	Phone        int                  `json:"phone,omitempty" bson:"phone,omitempty"`
	Owner        primitive.ObjectID   `json:"owner,omitempty" bson:"owner,omitempty"`
	Appointments []primitive.ObjectID `json:"appointments,omitempty" bson:"appointments,omitempty"`
	Employees    []primitive.ObjectID `json:"employees,omitempty" bson:"employees,omitempty"`
}

var (
	storeCollection = resources.Client.StoreCollection()
)

func (s Store) GetAll() ([]Store, *resterrors.RestError) {
	filter := bson.M{}
	cursor, err := getMultipleEntities(filter, storeCollection)
	if err != nil {
		return nil, err
	}

	list := make([]Store, 0)
	for cursor.Next(context.Background()) {
		var store Store
		if err := cursor.Decode(&store); err != nil {
			return nil, resterrors.NewInternalServerError("error parsing entity")
		}
		list = append(list, store)
	}
	return list, nil
}

func (s *Store) Save() (*primitive.ObjectID, *resterrors.RestError) {
	return saveEntity(&s, storeCollection)
}

func (s *Store) GetById() *resterrors.RestError {
	filter := bson.M{"id": s.ID}
	return getEntity(&s, filter, storeCollection)
}

func (s Store) Find(searchTerm string) ([]Store, *resterrors.RestError) {
	filter := bson.M{"name": getIgnoreCaseSearch(searchTerm)}

	cursor, err := getMultipleEntities(filter, storeCollection)
	if err != nil {
		return nil, err
	}

	list := make([]Store, 0)
	for cursor.Next(context.Background()) {
		var store Store
		if err := cursor.Decode(&store); err != nil {
			return nil, resterrors.NewInternalServerError("error parsing entity")
		}

		list = append(list, store)
	}
	return list, nil

}

func (s *Store) Update() *resterrors.RestError {
	filter := bson.M{"_id": s.ID}

	updating := bson.M{
		"$set": bson.M{
			"name":         s.Name,
			"address":      s.Address,
			"phone":        s.Phone,
			"owner":        s.Owner,
			"appointments": s.Appointments,
			"employees":    s.Employees,
		},
	}
	return updateEntity(&s, filter, updating, storeCollection)

}
