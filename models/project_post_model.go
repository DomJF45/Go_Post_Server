package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectData struct {
	Id      primitive.ObjectID `json:"id,omitempty"`
	Title   string             `json:"title,omitempty" validate:"required"`
	Content string             `json:"content,omitempty" validate:"required"`
	Img     string             `json:"img,omitempty"`
}

type ProjectPageData struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	ProjectName string             `json:"projectName,omitempty" validate:"required"`
	Posts       []ProjectData      `json:"posts,omitempty"`
	Stack       []string           `json:"stack,omitempty"`
}
