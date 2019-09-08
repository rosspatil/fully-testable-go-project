package mongo

import (
	"fully-testable-go-project/internal/models"

	"gopkg.in/mgo.v2/bson"
)

func (mg MongoClient) SaveStudent(student models.Student) error {
	_, err := mg.SaveData(models.StudentCollection, student)
	return err
}

func (mg MongoClient) GetStudent(name string) (models.Student, error) {
	student := models.Student{}
	err := mg.GetData(models.StudentCollection, bson.M{"name": name}, nil, &student)
	return student, err
}

func (mg MongoClient) UpdateStudent(name string, data bson.M) error {
	return mg.Update(models.StudentCollection, bson.M{"name": name}, data)
}

func (mg MongoClient) DeleteStudent(name string) error {
	return mg.DeleteData(models.StudentCollection, bson.M{"name": name})

}
