package models

import (
	"context"

	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type Service struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"`
	Category    string             `json:"category,omitempty" bson:"category,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
}

var (
	serviceCollection = resources.Client.ServiceCollection()
)

func (s Service) GetAll() ([]Service, *resterrors.RestError) {
	filter := bson.M{}
	cursor, err := getMultipleEntities(filter, serviceCollection)
	if err != nil {
		return nil, err
	}

	list := make([]Service, 0)
	for cursor.Next(context.Background()) {
		var service Service
		if err := cursor.Decode(&service); err != nil {
			return nil, resterrors.NewInternalServerError("error parsing entity")
		}
		list = append(list, service)
	}
	return list, nil
}

func (s *Service) Save() (*primitive.ObjectID, *resterrors.RestError) {
	return saveEntity(&s, serviceCollection)
}

func (s *Service) GetById() *resterrors.RestError {
	filter := bson.M{"_id": s.ID}
	return getEntity(&s, filter, serviceCollection)
}

func (s Service) Find(searchTerm string) ([]Service, *resterrors.RestError) {
	filter := bson.M{"name": getIgnoreCaseSearch(searchTerm)}

	cursor, err := getMultipleEntities(filter, serviceCollection)
	if err != nil {
		return nil, err
	}

	list := make([]Service, 0)
	for cursor.Next(context.Background()) {
		var service Service
		if err := cursor.Decode(&service); err != nil {
			return nil, resterrors.NewInternalServerError("error parsing entity")
		}

		list = append(list, service)
	}
	return list, nil

}

func (s *Service) Update() *resterrors.RestError {
	filter := bson.M{"_id": s.ID}

	updating := bson.M{
		"$set": bson.M{
			"name":  s.Name,
			"price": s.Price,
		},
	}
	return updateEntity(&s, filter, updating, serviceCollection)

}
