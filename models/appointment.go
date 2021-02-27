package models

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Appointment struct for appointment
// hold storeID since it a on-to-millions relationship (1 store -> m appointments)
type Appointment struct {
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	TechinicianID primitive.ObjectID   `json:"technicianId,omitempty" bson:"technicianId,omitempty"`
	CustomerID    primitive.ObjectID   `json:"customerId,omitempty" bson:"customerId,omitempty"`
	Time          time.Time            `json:"time,omitempty" bson:"time,omitempty"`
	Services      []primitive.ObjectID `json:"services,omitempty" bson:"services,omitempty"`
	StoreID       primitive.ObjectID   `json:"storeId,omitempty" bson:"storeId,omitempty"`
}

var (
	appCollection = resources.Client.AppointmentCollection()
)

// Save save appointment
func (a *Appointment) Save() (*primitive.ObjectID, *resterrors.RestError) {
	return saveEntity(&a, appCollection)
}

// GetByID get appointment by Id
func (a *Appointment) GetByID() *resterrors.RestError {
	filter := bson.M{"_id": a.ID}
	return getEntity(a, filter, appCollection)
}

// Update  appointment
func (a *Appointment) Update() *resterrors.RestError {
	filter := bson.M{"_id": a.ID}

	updating := bson.M{
		"$set": bson.M{
			"technicianId": a.TechinicianID,
			"customerId":   a.CustomerID,
			"time":         a.Time,
			"services":     a.Services,
		}}

	return updateEntity(a, filter, updating, appCollection)
}

//Delete appointment
func (a *Appointment) Delete() *resterrors.RestError {
	filter := bson.M{"_id": a.ID}
	return deleteEntity(&a, filter, appCollection)
}

//Find appointment
func (a *Appointment) Find(search string) ([]Appointment, *resterrors.RestError) {
	filter := bson.M{
		"$or": []interface{}{
			bson.M{
				"technicianId": getIgnoreCaseSearch(search),
				"customerId":   getIgnoreCaseSearch(search),
				"time":         getIgnoreCaseSearch(search),
				"services":     getIgnoreCaseSearch(search),
			},
		},
	}

	cursor, err := getMultipleEntities(filter, appCollection)
	if err != nil {
		return nil, err
	}
	list := make([]Appointment, 0)

	for cursor.Next(context.Background()) {
		var appointment Appointment
		if err := cursor.Decode(&appointment); err != nil {
			return nil, resterrors.NewInternalServerError("error parsing entity")
		}
		list = append(list, appointment)
	}
	return list, nil
}

// GetAll appointment
func (a *Appointment) GetAll() ([]Appointment, *resterrors.RestError) {
	filter := bson.M{}

	cursor, err := getMultipleEntities(filter, appCollection)
	if err != nil {
		return nil, err
	}
	list := make([]Appointment, 0)
	for cursor.Next(context.Background()) {
		var appointment Appointment
		if err := cursor.Decode(&appointment); err != nil {
			return nil, resterrors.NewBadRequestError("error parsing Appointment")
		}
		list = append(list, appointment)
	}
	return list, nil
}
func (a *Appointment) FindInTimeRange(fromDate, toDate time.Time) ([]Appointment, *resterrors.RestError) {
	filter := bson.M{
		"time": bson.M{
			"$gt": fromDate,
			"$lt": toDate,
		},
	}
	cursor, err := getMultipleEntities(filter, appCollection)
	if err != nil {
		return nil, err
	}
	list := make([]Appointment, 0)
	for cursor.Next(context.Background()) {
		var appointment Appointment
		if err := cursor.Decode(&appointment); err != nil {
			return nil, resterrors.NewBadRequestError("error parsing Appointment")
		}
		list = append(list, appointment)
	}
	return list, nil
}

// Validate an appointment
func (a Appointment) Validate() *resterrors.RestError {
	if err := validation.ValidateStruct(
		&a,
		validation.Field(&a.TechinicianID, validation.Required, validation.NotNil),
		validation.Field(&a.CustomerID, validation.Required, validation.NotNil),
		validation.Field(&a.StoreID, validation.Required, validation.NotNil),
		validation.Field(&a.Time, validation.Required, validation.NotNil),
		// validation.Field(&a.Services, validation.Required, validation.NotNil),
	); err != nil {
		return resterrors.NewBadRequestError(fmt.Sprintf("Invalid attribute %s!", err.Error()))
	}
	return nil
}
