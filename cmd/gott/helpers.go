package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func unwrapObjectID(oid primitive.ObjectID) string {
	return oid.String()[10:34]
}

func wrapObjectID(oid string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(oid)
}
