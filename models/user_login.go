package models

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/samderlust/spa_manager/utils/resterrors"
)

type UserLogin struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

func (u UserLogin) Validate() *resterrors.RestError {
	if err := validation.ValidateStruct(
		&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Match(regexp.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`))),
	); err != nil {
		return resterrors.NewBadRequestError(fmt.Sprintf("Invalid Info: %s", err.Error()))
	}
	return nil
}
