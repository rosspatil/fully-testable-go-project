package models

type Student struct {
	Name         string `json:"name" bson:"name"`
	CreationDate int64  `json:"creationDate" bson:"creationDate"`
	DOB          int64  `json:"dob" bson:"dob"`
}

const (
	StudentCollection = "student"
)
