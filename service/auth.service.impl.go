package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/models"
	"github.com/s1Sharp/s1-tts-restapi/internal/storage"
	"github.com/s1Sharp/s1-tts-restapi/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceImpl struct {
	mongoStorage *storage.MongoStorage
	ctx          context.Context
	config       *config.Config
}

func NewAuthService(storage *storage.MongoStorage, ctx context.Context, config *config.Config) AuthService {
	return &AuthServiceImpl{storage, ctx, config}
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBUserResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.LastAttemptAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = uc.config.AllVerifiedByDefault
	user.Role = "user"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	res, err := uc.mongoStorage.UserCollection().InsertOne(uc.ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, err
	}

	// Create a unique index for the email field
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := uc.mongoStorage.UserCollection().Indexes().CreateOne(uc.ctx, index); err != nil {
		return nil, errors.New("could not create index for email")
	}

	var newUser *models.DBUserResponse
	query := bson.M{"_id": res.InsertedID}

	err = uc.mongoStorage.UserCollection().FindOne(uc.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *AuthServiceImpl) SignInUser(*models.SignInInput) (*models.DBUserResponse, error) {
	return nil, nil
}
