package resources

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//Client export mongo client
var (
	Client bookingCoreInterface = &bookingCoreClient{}
	client *mongo.Client
)

type bookingCoreInterface interface {
	TechnicianCollection() *mongo.Collection
	AvailabilityCollection() *mongo.Collection
	CustomerCollection() *mongo.Collection
	AppointmentCollection() *mongo.Collection
	ServiceCollection() *mongo.Collection
	UserCollection() *mongo.Collection
}

type bookingCoreClient struct{}

func init() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
}

func (b *bookingCoreClient) UserCollection() *mongo.Collection {
	c := client.Database("booking_core").Collection("users")
	index := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	if _, err := c.Indexes().CreateOne(
		context.Background(),
		index,
	); err != nil {
		log.Fatal("failed create Index")
	}
	return c
}

//TechnicianCollection get technicians collection
func (b *bookingCoreClient) TechnicianCollection() *mongo.Collection {
	c := client.Database("booking_core").Collection("technicians")
	index := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	if _, err := c.Indexes().CreateOne(
		context.Background(),
		index,
	); err != nil {
		log.Fatal("failed create Index")
	}
	return c
}

//AvailabilityCollection get availabilities collection
func (b *bookingCoreClient) AvailabilityCollection() *mongo.Collection {
	return client.Database("booking_core").Collection("availabilities")
}

//CustomerCollection get customers collection
func (b *bookingCoreClient) CustomerCollection() *mongo.Collection {
	c := client.Database("booking_core").Collection("customers")
	index := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	if _, err := c.Indexes().CreateOne(
		context.Background(),
		index,
	); err != nil {
		log.Fatal("failed create Index")
	}
	return c
}

//AppointmentCollection get appointment collection
func (b *bookingCoreClient) AppointmentCollection() *mongo.Collection {
	return client.Database("booking_core").Collection("appointments")
}

// ServiceCollection get services collection
func (b *bookingCoreClient) ServiceCollection() *mongo.Collection {
	return client.Database("booking_core").Collection("services")
}
