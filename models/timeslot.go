package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TimeSlot struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Professor primitive.ObjectID `bson:"professor" json:"professor"`
	Slot      string             `bson:"slot" json:"slot"`
	Booked    bool               `bson:"booked" json:"booked"`
}
