package models

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Technician struct {
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Services       []primitive.ObjectID `json:"services,omitempty" bson:"services,omitempty"`
	Availabilities []primitive.ObjectID `json:"availabilities,omitempty" bson:"availabilities,omitempty"`
	UserID         primitive.ObjectID   `json:"userId,omitempty" bson:"userId,omitempty"`
}

type PopulatedTechnician struct {
	Technician
	_id primitive.ObjectID `json:"-" bson:"-"`
	User
}

var (
	techCollection = resources.Client.TechnicianCollection()
)

func (t *Technician) Marshall() (*PopulatedTechnician, *resterrors.RestError) {
	var user User
	user.ID = t.UserID
	if err := user.FindByID(); err != nil {
		return nil, err
	}
	tech := new(PopulatedTechnician)
	tech.Technician = *t
	mUser := user.Marshall()
	tech.User = *mUser
	return tech, nil
}

func (t *Technician) GetAll() ([]Technician, *resterrors.RestError) {
	filter := bson.M{}
	cursor, err := getMultipleEntities(filter, techCollection)
	if err != nil {
		return nil, err
	}

	list := make([]Technician, 0)

	for cursor.Next(context.Background()) {
		var technician Technician
		if err := cursor.Decode(&technician); err != nil {
			return nil, resterrors.NewInternalServerError("error parsing entity")
		}
		list = append(list, technician)
	}
	return list, nil
}

func (t *Technician) Save() (*primitive.ObjectID, *resterrors.RestError) {
	return saveEntity(&t, techCollection)
}

func (t *Technician) GetByID() *resterrors.RestError {
	filter := bson.M{"_id": t.ID}
	return getEntity(&t, filter, techCollection)
}

func (t *Technician) Find(search string) ([]Technician, *resterrors.RestError) {
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"firstName": getIgnoreCaseSearch(search)},
			bson.M{"lastName": getIgnoreCaseSearch(search)},
			bson.M{"email": getIgnoreCaseSearch(search)},
			bson.M{"phoneNumber": getIgnoreCaseSearch(search)},
		},
	}

	cursor, err := getMultipleEntities(filter, techCollection)
	if err != nil {
		return nil, err
	}
	list := make([]Technician, 0)

	for cursor.Next(context.Background()) {
		var technician Technician
		if err := cursor.Decode(&technician); err != nil {
			return nil, resterrors.NewInternalServerError("error parsing entity")
		}
		list = append(list, technician)
	}
	return list, nil
}

func (t *Technician) Update() *resterrors.RestError {

	return updateEntity(
		&t,
		bson.M{"_id": t.ID},
		bson.M{
			"$set": bson.M{}},
		techCollection,
	)
}

func (t *Technician) Delete() *resterrors.RestError {
	return deleteEntity(&t, bson.M{"_id": t.ID}, techCollection)
}

//Validate validate required fields of tecnician
func (t Technician) Validate() *resterrors.RestError {
	if err := validation.ValidateStruct(
		&t,
	); err != nil {
		return resterrors.NewBadRequestError(fmt.Sprintf("Invalid attribute %s!", err.Error()))
	}
	return nil
}
