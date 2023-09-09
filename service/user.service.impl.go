package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/s1Sharp/s1-tts-restapi/internal/models"
	"github.com/s1Sharp/s1-tts-restapi/internal/storage"
	"github.com/s1Sharp/s1-tts-restapi/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserServiceImpl struct {
	mongoStorage *storage.MongoStorage
	ctx          context.Context
}

func NewUserServiceImpl(storage *storage.MongoStorage, ctx context.Context) UserService {
	return &UserServiceImpl{storage, ctx}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse

	query := bson.M{"_id": oid}
	err := us.mongoStorage.UserCollection().FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse

	query := bson.M{"email": strings.ToLower(email)}
	log.Warnf("find user %s ", email)
	err := us.mongoStorage.UserCollection().FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (uc *UserServiceImpl) UpdateUserById(id string, data *models.UpdateInput) (*models.DBResponse, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return &models.DBResponse{}, err
	}

	fmt.Println(data)

	obId, _ := primitive.ObjectIDFromHex(id)

	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: doc}}
	result := uc.mongoStorage.UserCollection().FindOneAndUpdate(uc.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedUser *models.DBResponse
	if err := result.Decode(&updatedUser); err != nil {
		return nil, errors.New("no document with that id exists")
	}

	return updatedUser, nil
}

func (uc *UserServiceImpl) RemoveUserById(id string) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: obId}}

	result, err := uc.mongoStorage.UserCollection().DeleteOne(uc.ctx, query)
	if err != nil {
		log.Printf("Remove User By id fail to delete one, %s", err)
		return err
	}
	if result.DeletedCount == 0 {
		log.Printf("Remove User By count of deleted is 0, _id=%s", obId)
	}
	return nil
}
