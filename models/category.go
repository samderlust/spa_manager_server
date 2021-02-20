package models

import (
	"context"

	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type Category struct {
	ID   primitive.ObjectID `json:"id, omitempty" bson:"_id,omitempty"`
	name string             `json:"name, omitempty" bson:"name,omitempty"`
}

var (
	categoryCollection = resources.Client.CategoryCollection()
)

func (c *Category) GetAll() ([]Category, *resterrors.RestError) {
	filter := bson.M{}
	cursor, err := getMultipleEntities(filter, categoryCollection)
	if err != nil {
		return nil, err
	}

	list := make([]Category, 0)
	for cursor.Next(context.Background()) {
		var category Category
		if err := cursor.Decode(&category); err != nil {
			return nil, resterrors.NewBadRequestError("error parsing Category")
		}
		list = append(list, category)
	}
	return list, nil
}

func (c *Category) Save() (*primitive.ObjectID, *resterrors.RestError) {
	return saveEntity(&c, categoryCollection)
}

func (c *Category) GetByID() *resterrors.RestError {
	filter := bson.M{"_id": c.ID}
	return getEntity(c, filter, categoryCollection)
}
