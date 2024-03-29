package models

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/logger"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Store struct
type Store struct {
	ID        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string               `json:"name,omitempty" bson:"name,omitempty"`
	Address   string               `json:"address,omitempty" bson:"address,omitempty"`
	Phone     int                  `json:"phone,omitempty" bson:"phone,omitempty"`
	Owner     primitive.ObjectID   `json:"owner,omitempty" bson:"owner,omitempty"`
	Services  []primitive.ObjectID `json:"services," bson:"services"`
	Employees []primitive.ObjectID `json:"employees" bson:"employees"`
}

// PopulatedStore is used to populated owner and employees
type PopulatedStore struct {
	Store
	Owner        User          `json:"owner" bson:"owner"`
	Employees    []Technician  `json:"employees," bson:"employees,"`
	Appointments []Appointment `json:"appointments" bson:"appointments"`
	Services     []Service     `json:"services," bson:"services,"`
}

type address struct {
}

var (
	storeCollection = resources.Client.StoreCollection()
)

// Puplate the store
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

	eList := make([]Technician, 0)
	for _, e := range s.Employees {
		var tech Technician
		tech.ID = e

		if err := tech.GetByID(); err != nil {
			logger.Info(err.Message)
			// return nil
		}
		eList = append(eList, tech)
	}
	pStore.Employees = eList

	sList := make([]Service, 0)
	for _, e := range s.Services {
		var ser Service
		ser.ID = e

		if err := ser.GetById(); err != nil {
			logger.Info(err.Message)
		}
		sList = append(sList, ser)
	}
	pStore.Services = sList

	return pStore
}

// GetAll stores
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

// Save store
func (s *Store) Save() (*primitive.ObjectID, *resterrors.RestError) {
	return saveEntity(&s, storeCollection)
}

// GetByID Store
func (s *Store) GetByID() *resterrors.RestError {
	filter := bson.M{"id": s.ID}
	return getEntity(&s, filter, storeCollection)
}

// FindByOwnerID get store by owner Id
func (s *Store) FindByOwnerID() *resterrors.RestError {
	filter := bson.M{"owner": s.Owner}
	return getEntity(&s, filter, storeCollection)
}

// Find store
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

// Update store
func (s *Store) Update() *resterrors.RestError {
	filter := bson.M{"_id": s.ID}
	updating := bson.M{
		"$set": bson.M{
			"name":      s.Name,
			"address":   s.Address,
			"phone":     s.Phone,
			"owner":     s.Owner,
			"employees": s.Employees,
		},
	}
	return updateEntity(&s, filter, updating, storeCollection)

}

// AddService into store
func (s *Store) AddService(service Service) *resterrors.RestError {
	filter := bson.M{"_id": s.ID}

	updating := bson.M{
		"$push": bson.M{
			"services": service.ID,
		},
	}

	return updateEntity(&s, filter, updating, storeCollection)
}

// AddTechnician into store
func (s *Store) AddTechnician(tech Technician) *resterrors.RestError {
	filter := bson.M{"_id": s.ID}
	updating := bson.M{
		"$push": bson.M{
			"employees": tech.ID,
		},
	}
	return updateEntity(&s, filter, updating, storeCollection)
}

// Validate store
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
