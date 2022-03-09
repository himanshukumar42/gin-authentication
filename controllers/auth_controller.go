package controllers

import (
	"context"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/himanshuk42/gin-authentication/configs"
	"github.com/himanshuk42/gin-authentication/models"
	"github.com/himanshuk42/gin-authentication/responses"
	"github.com/himanshuk42/gin-authentication/utils/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var authCollection *mongo.Collection = configs.GetCollection(configs.DB, "autho")
var authValidate = validator.New()

func beforeSave(u *models.RegisterInput) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(ctx context.Context, username string, password string) (string, error) {
	var validateUser models.RegisterInput
	userErr := authCollection.FindOne(ctx, bson.M{"username": username}).Decode(&validateUser)
	if userErr != nil {
		return "", userErr
	}
	passErr := verifyPassword(password, validateUser.Password)
	if passErr != nil && passErr == bcrypt.ErrMismatchedHashAndPassword {
		return "", passErr
	}

	token, err := token.GenerateToken(validateUser.Id.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var input models.RegisterInput

		defer cancel()

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if validationErr := authValidate.Struct(&input); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}
		inputUser := models.RegisterInput{
			Id:       primitive.NewObjectID(),
			Username: input.Username,
			Password: input.Password,
		}

		beforeSave(&inputUser)

		_, err := authCollection.InsertOne(ctx, inputUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": inputUser}})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var input models.RegisterInput

		defer cancel()

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		token, err := LoginCheck(ctx, input.Username, input.Password)
		if err != nil {
			c.JSON(http.StatusForbidden, responses.UserResponse{Status: http.StatusForbidden, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"access_token": token}})
	}
}
