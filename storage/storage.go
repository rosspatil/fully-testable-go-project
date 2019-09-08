package storage

import (
	"fully-testable-go-project/internal/models"
	"gopkg.in/mgo.v2/bson"

)

var storageInstance Storage

type Storage interface {
	GetStudent(string) (models.Student, error)
	SaveStudent(models.Student) error
	DeleteStudent(name string)error
	UpdateStudent(name string, data bson.M) error
}

func Init(dbInstance Storage) {
	storageInstance = dbInstance
}

func GetStorageInstance() Storage {
	return storageInstance
}
