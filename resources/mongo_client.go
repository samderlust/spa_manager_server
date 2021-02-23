package resources

import (
	"context"
	"log"
	"time"

	"github.com/samderlust/spa_manager/utils/logger"
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
	CategoryCollection() *mongo.Collection
}

type bookingCoreClient struct{}

func init() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var err error
	uri := "mongodb+srv://samderlust:P%40ssword1@cluster0.l4grx.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Info(err.Error())
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

func (b *bookingCoreClient) CategoryCollection() *mongo.Collection {
	c := client.Database("booking_core").Collection("categories")
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
