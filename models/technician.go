package models

import (
	"context"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Technician struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName   string             `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName    string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	PhoneNumber string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	// Availabilities []bson.ObjectID `json:"availabilities,omitempty" bson:"availabilities,omitempty"`
	// Bookings       []bson.ObjectID `json:"bookings,omitempty" bson:"bookings,omitempty"`
}

var (
	techCollection = resources.Client.TechnicianCollection()
)

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
			"$set": bson.M{
				"firstName":   t.FirstName,
				"lastName":    t.LastName,
				"phoneNumber": t.PhoneNumber,
				"email":       t.Email,
			}},
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
		validation.Field(&t.FirstName, validation.Required, validation.NotNil, validation.Length(1, 10)),
		validation.Field(&t.LastName, validation.Required, validation.Length(1, 10)),
		validation.Field(&t.PhoneNumber, validation.Required, validation.NotNil, validation.Match(regexp.MustCompile(`^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`))),
		validation.Field(&t.Email, validation.Match(regexp.MustCompile(`^([\w\.\-]+)@([\w\-]+)((\.(\w){2,3})+)$`))),
	); err != nil {
		return resterrors.NewBadRequestError(fmt.Sprintf("Invalid attribute %s!", err.Error()))
	}
	return nil
}
