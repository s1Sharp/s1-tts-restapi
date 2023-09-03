package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/s1Sharp/s1-tts-restapi/internal/models"
	"github.com/s1Sharp/s1-tts-restapi/internal/storage"
	"github.com/s1Sharp/s1-tts-restapi/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskServiceImpl struct {
	mongoStorage *storage.MongoStorage
	ctx          context.Context
}

func NewTaskService(mongoStorage *storage.MongoStorage, ctx context.Context) TaskService {
	return &TaskServiceImpl{mongoStorage, ctx}
}

func (p *TaskServiceImpl) CreateTask(task *models.CreateTaskScheme) (*models.DBTaskScheme, error) {

	var insertTask models.InsertTaskScheme
	// TODO self copying
	insertTask.CreatedAt = time.Now()
	insertTask.UpdatedAt = insertTask.CreatedAt
	insertTask.Completed = false
	insertTask.DownloadUrl = ""
	insertTask.TaskUuid = uuid.New()
	// TODO must be a login
	insertTask.User = task.User
	insertTask.Text = task.Text
	insertTask.IsSSML = task.IsSSML

	res, err := p.mongoStorage.TaskCollection().InsertOne(p.ctx, &insertTask)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("task with that title already exists")
		}
		return nil, err
	}

	//opt := options.Index()
	//opt.SetUnique(true)
	//
	//index := mongo.IndexModel{Keys: bson.M{"title": 1}, Options: opt}
	//
	//if _, err := p.mongoStorage.TaskCollection().Indexes().CreateOne(p.ctx, index); err != nil {
	//	return nil, errors.New("could not create index for title")
	//}

	var newTask *models.DBTaskScheme
	query := bson.M{"_id": res.InsertedID}
	if err = p.mongoStorage.TaskCollection().FindOne(p.ctx, query).Decode(&newTask); err != nil {
		return nil, err
	}

	return newTask, nil
}

func (p *TaskServiceImpl) updateTask(id string, data *models.DBTaskScheme) (*models.DBTaskScheme, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return nil, err
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := p.mongoStorage.TaskCollection().FindOneAndUpdate(p.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedPost *models.DBTaskScheme

	if err := res.Decode(&updatedPost); err != nil {
		return nil, errors.New("no task with that Id exists")
	}

	return updatedPost, nil
}

func (p *TaskServiceImpl) TaskById(id string) (*models.DBTaskScheme, error) {
	obId, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": obId}

	var tasks *models.DBTaskScheme

	if err := p.mongoStorage.TaskCollection().FindOne(p.ctx, query).Decode(&tasks); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return tasks, nil
}

func (p *TaskServiceImpl) GetUserTasks(user string, page int, limit int) ([]*models.DBTaskScheme, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	skip := (page - 1) * limit

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	query := bson.M{"user": user}

	cursor, err := p.mongoStorage.TaskCollection().Find(p.ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	defer func() {
		e := cursor.Close(p.ctx)
		if e != nil {
			log.Fatalf("unexpected cursor close: %s", e)
		}
	}()

	var tasks []*models.DBTaskScheme

	for cursor.Next(p.ctx) {
		post := &models.DBTaskScheme{}
		err := cursor.Decode(post)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return []*models.DBTaskScheme{}, nil
	}

	return tasks, nil
}
