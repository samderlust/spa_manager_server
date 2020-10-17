package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/samderlust/spa_manager/models"
	"github.com/samderlust/spa_manager/utils/httputils"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("klsajflkdasdsfdsfsdjf")

var AuthHandler authHandlerI = &authHandler{}

type authHandler struct{}
type authHandlerI interface {
	SignUp(*fiber.Ctx) error
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

	return httputils.JSONCreatedResponse(c, user)
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
	validToken, err := _generateJWT()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	return httputils.JSONSuccessResponse(c, &fiber.Map{
		"user":  user.Marshall(),
		"token": validToken,
	})
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

func _generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "samderlust"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Printf("something went wrong : %s \n", err.Error())
		return "", err
	}
	return tokenString, nil
}
