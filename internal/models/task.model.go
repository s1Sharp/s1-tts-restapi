package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewTask(user string, text string, voiceType string, isSSML bool, taskUuid uuid.UUID) *CreateTaskScheme {
	return &CreateTaskScheme{User: user, Text: text, VoiceType: voiceType, IsSSML: isSSML}
}

// TODO create restriction for fields

type CreateTaskScheme struct {
	User      string `json:"user" bson:"user" binding:"required"`
	Text      string `json:"text" bson:"text" binding:"required"`
	VoiceType string `json:"voice_type,omitempty" bson:"voice_type"`
	IsSSML    bool   `json:"is_ssml,omitempty" bson:"is_ssml"`
}

type InsertTaskScheme struct {
	// CreateTaskScheme
	User      string `bson:"user"`
	Text      string `bson:"text"`
	VoiceType string `bson:"voice_type"`
	IsSSML    bool   `bson:"is_ssml"`

	// new fields
	TaskUuid    uuid.UUID `bson:"task_uuid"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	Completed   bool      `bson:"completed"`
	DownloadUrl string    `bson:"download_url"`
}

type DBTaskScheme struct {
	// CreateTaskScheme
	User      string    `bson:"user"`
	Text      string    `bson:"text"`
	VoiceType string    `bson:"voice_type"`
	IsSSML    bool      `bson:"is_ssml"`
	TaskUuid  uuid.UUID `bson:"task_uuid"`
	// InsertTaskScheme
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	Completed   bool      `bson:"completed"`
	DownloadUrl string    `bson:"download_url"`
	// new fields
	ID primitive.ObjectID `bson:"_id"`
}

type GetTaskSchemeResponse struct {
	TaskUuid    uuid.UUID          `bson:"task_uuid" json:"task_uuid,omitempty"`
	Completed   bool               `bson:"completed" json:"completed,omitempty"`
	DownloadUrl string             `bson:"download_url" json:"download_url,omitempty"`
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

func DbTaskToTaskResponse(t []*DBTaskScheme) []*GetTaskSchemeResponse {
	var tasks []*GetTaskSchemeResponse
	for _, task := range t {
		tasks = append(tasks, task.ToTaskResponse())
	}
	return tasks
}

func (t *DBTaskScheme) ToTaskResponse() *GetTaskSchemeResponse {
	return &GetTaskSchemeResponse{
		TaskUuid:    t.TaskUuid,
		Completed:   t.Completed,
		DownloadUrl: t.DownloadUrl,
		ID:          t.ID,
	}
}
