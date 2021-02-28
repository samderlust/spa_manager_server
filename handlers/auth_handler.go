package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/bookstore_utils-go/resterrors"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("klsajflkdasdsfdsfsdjf")

// AuthHandler export AuthHandler
var AuthHandler authHandlerI = &authHandler{}

type authHandler struct{}
type authHandlerI interface {
	SignUp(*fiber.Ctx) error
	SignIn(*fiber.Ctx) error
	ChangePassword(*fiber.Ctx) error
}

//UserChangePassword type of user that want to change password
type UserChangePassword struct {
	models.UserLogin
	NewPassword string `json:"newPassword,omitempty" bson:"newPassword,omitempty"`
}

func (h authHandler) SignUp(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}

	if err := user.Validate(); err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	hashedPwd := _hashAndSalt([]byte(user.Password))
	user.Password = hashedPwd

	ID, err := user.Save()
	if err != nil {
		return httputils.JSONResponseModelError(c, err)
	}
	user.ID = *ID
	user = user.Marshall()

	validToken, serr := _generateJWT(user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(serr)
	}

	user.Token = validToken

	return httputils.JSONSuccessResponse(c, user.Marshall())
}

func (h authHandler) SignIn(c *fiber.Ctx) error {
	userLogin := new(models.UserLogin)

	if err := c.BodyParser(userLogin); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}

	if err := userLogin.Validate(); err != nil {
		return httputils.JSONResponseModelError(c, err)
	}

	user := new(models.User)
	if err := user.FindByEmail(userLogin.Email); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid Email or Password"})
	}

	if ok := _comparePasswords(user.Password, userLogin.Password); !ok {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid Email or Password"})
	}
	validToken, err := _generateJWT(user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	user.Token = validToken

	return httputils.JSONSuccessResponse(c, user.Marshall())
}

func (h authHandler) ChangePassword(c *fiber.Ctx) error {
	userChangePassword := new(UserChangePassword)

	if err := c.BodyParser(userChangePassword); err != nil {
		return httputils.JSONParamInvalidResponse(c, err)
	}
	currentUser := new(models.User)
	if err := currentUser.FindByEmail(userChangePassword.Email); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid Email or Password"})
	}

	if ok := _comparePasswords(currentUser.Password, userChangePassword.Password); !ok {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid Email or Password"})
	}
	hashedPwd := _hashAndSalt([]byte(userChangePassword.NewPassword))
	currentUser.Password = hashedPwd

	if err := currentUser.Update(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	validToken, err := _generateJWT(currentUser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	currentUser.Token = validToken

	return httputils.JSONSuccessResponse(c, currentUser.Marshall())
}

func _hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Print(err)
	}
	return string(hash)
}

func _comparePasswords(hashedPwd string, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		log.Print(err)
		return false
	}
	return true
}

func _generateJWT(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = user.ID
	claims["role"] = user.Role
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", resterrors.NewInternalServerError("Something wrong...", err)
	}
	return tokenString, nil
}
