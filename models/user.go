package models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/samderlust/spa_manager/resources"
	"github.com/samderlust/spa_manager/utils/resterrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = resources.Client.UserCollection()

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Role     string             `json:"role,omitempty" bson:"role,omitempty"`
	Token    string             `json:"token,omitempty"`
}

func (u *User) Marshall() *User {
	u.Password = ""
	return u
}

func (u *User) Save() (*primitive.ObjectID, *resterrors.RestError) {
	return saveEntity(&u, userCollection)
}

func (u *User) Update() *resterrors.RestError {
	filter := bson.M{"_id": u.ID}
	updating := bson.M{
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
	}
	return updateEntity(&u, filter, updating, userCollection)
}

func (u *User) Delete() *resterrors.RestError {
	filter := bson.M{"_id": u.ID}
	return deleteEntity(&u, filter, userCollection)
}

func (u *User) FindByEmail(email string) *resterrors.RestError {
	filter := bson.M{
		"email": email,
	}
	return getEntity(&u, filter, userCollection)
}

func (u User) Validate() *resterrors.RestError {
	if err := validation.ValidateStruct(
		&u,
		validation.Field(&u.Username, validation.Required, validation.NotNil, validation.Length(4, 18)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(4, 255)),
		validation.Field(&u.Role, validation.Required, validation.In("admin", "technician")),
	); err != nil {
		return resterrors.NewBadRequestError(fmt.Sprintf("Invalid Info: %s", err.Error()))
	}
	return nil
}
