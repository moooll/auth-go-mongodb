package main

import (
	"testing"
//	"context"
//	"reflect"
	"log"

	"github.com/gobuffalo/uuid"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
)


// func TestWriteUser(t *testing.T) {
// 	id := uuid.Must(uuid.NewV4())
// 	connectToTheDB()
// 	res := writeUser(id)
// 	var result bson.M
// 	if err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result); err != nil && err == mongo.ErrNoDocuments || res == nil {
// 		t.Log("result id type ", reflect.TypeOf(res))
// 		t.Errorf("user is not not found in the DB, id: %s, error: %s", res, err)
// 	}
// }

func TestFindUser(t *testing.T) {
	uid, err := uuid.FromString("c9016958-50ee-46c0-b242-129c0092d88a")
	if err != nil {
		log.Fatalf("error parsing uuid from string: %s", err)
	}
	isPresent, err := findUser(uid)
	if isPresent == false || err != nil {
		t.Error("test failed, user id not found")
	}
}

