package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Appointment struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Student      primitive.ObjectID `bson:"student" json:"student"`
	Professor    primitive.ObjectID `bson:"professor" json:"professor"`
	Availability primitive.ObjectID `bson:"availability" json:"availability"`
	Status       string             `bson:"status" json:"status"`
}
