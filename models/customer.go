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

//Customer describe customer entity
type Customer struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName       string             `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName        string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	PhoneNumber     string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	Email           string             `json:"email,omitempty" bson:"email,omitempty"`
	NextAppointment primitive.ObjectID `json:"nextApointment,omitempty" bson:"nextApointment,omitempty"`
}

var (
	cusCollection = resources.Client.CustomerCollection()
)

func (c *Customer) GetAll() ([]Customer, *resterrors.RestError) {
	filter := bson.M{}
	cursor, err := getMultipleEntities(filter, cusCollection)
	if err != nil {
		return nil, err
	}
	var list []Customer
	for cursor.Next(context.Background()) {
		var customer Customer
		if err := cursor.Decode(&customer); err != nil {
			return nil, resterrors.NewInternalServerError("error parsing Customer")
		}
		list = append(list, customer)
	}
	return list, nil
}

//Save save a customer into DB
func (c *Customer) Save() (*primitive.ObjectID, *resterrors.RestError) {
	return saveEntity(&c, cusCollection)
}

func (c *Customer) GetByID() *resterrors.RestError {
	filter := bson.M{"_id": c.ID}
	return getEntity(c, filter, cusCollection)
}

func (c *Customer) UpdateNextAppointment(updatingID primitive.ObjectID) *resterrors.RestError {

	filter := bson.M{"_id": c.ID}
	updating := bson.M{
		"$set": bson.M{
			"nextApointment": updatingID,
		},
	}

	return updateEntity(&c, filter, updating, cusCollection)
}

func (c *Customer) Update() *resterrors.RestError {

	filter := bson.M{"_id": c.ID}

	updating := bson.M{
		"$set": bson.M{
			"firstName":      c.FirstName,
			"lastName":       c.LastName,
			"phoneNumber":    c.PhoneNumber,
			"email":          c.Email,
			"nextApointment": c.NextAppointment,
		}}
	return updateEntity(c, filter, updating, cusCollection)
}

func (c *Customer) Delete() *resterrors.RestError {
	return deleteEntity(&c, bson.M{"_id": c.ID}, cusCollection)
}

func (c *Customer) Find(search string) ([]Customer, *resterrors.RestError) {
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"firstName": getIgnoreCaseSearch(search)},
			bson.M{"lastName": getIgnoreCaseSearch(search)},
			bson.M{"email": getIgnoreCaseSearch(search)},
			bson.M{"phoneNumber": getIgnoreCaseSearch(search)},
		},
	}

	cursor, err := getMultipleEntities(filter, cusCollection)
	if err != nil {
		return nil, err
	}
	var list []Customer

	for cursor.Next(context.Background()) {
		var customer Customer
		if err := cursor.Decode(&customer); err != nil {
			return nil, resterrors.NewInternalServerError("error parsing entity")
		}
		list = append(list, customer)
	}
	return list, nil
}

//Validate validating customer's fields
func (c Customer) Validate() *resterrors.RestError {
	if err := validation.ValidateStruct(
		&c,
		validation.Field(&c.FirstName, validation.Required, validation.NotNil, validation.Length(1, 15)),
		validation.Field(&c.LastName, validation.Required, validation.Length(1, 15)),
		validation.Field(&c.PhoneNumber, validation.Required, validation.NotNil, validation.Match(regexp.MustCompile(`^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`))),
		validation.Field(&c.Email, validation.Match(regexp.MustCompile(`^([\w\.\-]+)@([\w\-]+)((\.(\w){2,3})+)$`))),
	); err != nil {
		return resterrors.NewBadRequestError(fmt.Sprintf("Invalid attribute %s!", err.Error()))
	}
	return nil
}
