package models

import (
	"context"
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/logger"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string               `json:"name,omitempty" bson:"name,omitempty"`
	Address      string               `json:"address,omitempty" bson:"address,omitempty"`
	Phone        int                  `json:"phone,omitempty" bson:"phone,omitempty"`
	Owner        primitive.ObjectID   `json:"owner,omitempty" bson:"owner,omitempty"`
	Appointments []primitive.ObjectID `json:"appointments,omitempty" bson:"appointments,omitempty"`
	Employees    []primitive.ObjectID `json:"employees,omitempty" bson:"employees,omitempty"`
}

type PopulatedStore struct {
	Store
	Owner        User          `json:"owner,omitempty" bson:"owner,omitempty"`
	Appointments []Appointment `json:"appointments," bson:"appointments,"`
	Employees    []User        `json:"employees," bson:"employees,"`
}

type address struct {
}

var (
	storeCollection = resources.Client.StoreCollection()
)

func (s *Store) Puplate() *PopulatedStore {
	pStore := new(PopulatedStore)
	pStore.ID = s.ID
	pStore.Phone = s.Phone
	pStore.Address = s.Address
	pStore.Name = s.Name

	owner := User{}
	owner.ID = s.Owner
	if err := owner.FindByID(); err != nil {
		logger.Info(err.Message)
		// return nil
	}
	pStore.Owner = *owner.Marshall()

	appList := make([]Appointment, 0)
	for _, a := range s.Appointments {
		var appointment Appointment
		appointment.ID = a

		if err := appointment.GetByID(); err != nil {

			logger.Info(err.Message)
			// return nil
		}
		appList = append(appList, appointment)
	}
	pStore.Appointments = appList

	eList := make([]User, 0)
	for _, e := range s.Employees {
		var emp User
		emp.ID = e

		if err := emp.FindByID(); err != nil {
			logger.Info(err.Message)
			// return nil
		}
		eList = append(eList, emp)
	}
	logger.Info("app")
	logger.Info(strconv.Itoa(len(pStore.Appointments)))
	pStore.Employees = eList

	return pStore
}

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

func (s Store) Validate() *resterrors.RestError {
	if err := validation.ValidateStruct(
		&s,
		validation.Field(&s.Name, validation.Required, validation.NotNil, validation.Length(4, 18)),
		validation.Field(&s.Address, validation.Required, validation.NotNil, validation.Length(4, 255)),
		validation.Field(&s.Phone, validation.Required, validation.NotNil),
	); err != nil {
		return resterrors.NewBadRequestError(fmt.Sprintf("Invalid Info: %s", err.Error()))
	}
	return nil
}
