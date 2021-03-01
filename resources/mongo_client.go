package resources

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/samderlust/spa_manager/utils/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//Client export mongo client
var (
	Client    bookingCoreInterface = &bookingCoreClient{}
	MgoClient *mongo.Client
)

type bookingCoreInterface interface {
	TechnicianCollection() *mongo.Collection
	AvailabilityCollection() *mongo.Collection
	CustomerCollection() *mongo.Collection
	AppointmentCollection() *mongo.Collection
	ServiceCollection() *mongo.Collection
	UserCollection() *mongo.Collection
	CategoryCollection() *mongo.Collection
	StoreCollection() *mongo.Collection
}

type bookingCoreClient struct{}

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	uri := os.Getenv("DB_URL")
	MgoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Info(err.Error())
		panic(err)
	}
}

func (b *bookingCoreClient) UserCollection() *mongo.Collection {
	c := MgoClient.Database("booking_core").Collection("users")
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

func (b *bookingCoreClient) CategoryCollection() *mongo.Collection {
	c := MgoClient.Database("booking_core").Collection("categories")
	index := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
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
	c := MgoClient.Database("booking_core").Collection("technicians")

	return c
}

//AvailabilityCollection get availabilities collection
func (b *bookingCoreClient) AvailabilityCollection() *mongo.Collection {
	return MgoClient.Database("booking_core").Collection("availabilities")
}

//StoreCollection get availabilities collection
func (b *bookingCoreClient) StoreCollection() *mongo.Collection {
	return MgoClient.Database("booking_core").Collection("stores")
}

//CustomerCollection get customers collection
func (b *bookingCoreClient) CustomerCollection() *mongo.Collection {
	c := MgoClient.Database("booking_core").Collection("customers")
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
	return MgoClient.Database("booking_core").Collection("appointments")
}

// ServiceCollection get services collection
func (b *bookingCoreClient) ServiceCollection() *mongo.Collection {
	return MgoClient.Database("booking_core").Collection("services")
}
