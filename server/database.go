package main

import (
	"log"
	"context"
	//"time"
	
	"github.com/gobuffalo/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"golang.org/x/crypto/bcrypt"
)
//User struct describes a user 
type User struct {
	ID uuid.UUID `json: "_id"`
	Session Session `json: "session_id"`
}

var client mongo.Client
const dbURI =  "mongodb://localhost:30001,localhost:30002,localhost:30003/?replicaSet=rs0"
var collection 	*mongo.Collection
//var cntx, _ = context.WithTimeout(context.TODO(), 10*time.Second)
//defer cancel()

func connectToTheDB() {
	//contxt, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	client, err  := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI)) 
	if err != nil {
		log.Fatal("problems connecting to the db:(",err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("authGo").Collection("users")
}

//check if user is in dB passing id as arg
func findUser(id uuid.UUID) (isPresent bool){
	//if user is present in DB isPresent is true
	filter := bson.D{{"_id", id}}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal("error finding users: ", err)
	}else {
		isPresent = true
	}
	return isPresent
}

func readUsers() (results User) {
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalf("error finding users: %s", err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()){
		err = cur.Decode(results)
		if err != nil {
			log.Fatalf("error decoding users: %s", err)
		}
	}
	return results
}

//writes new session to the DB
//edited 17:38
func insertNewSession(session Session) {
	opts := options.FindOneAndUpdate().SetUpsert(false)
	filter := bson.M{"_id":session.UserID}
	update := bson.M{"$set": bson.M{"_id": session.UserID, "session": bson.M{"session_id": session.ID, "refresh": session.Refresh, "expires_at":session.ExpiresAt, "is_session_over": session.IsSessionOver}}}
	err := collection.FindOneAndUpdate(context.Background(), filter, update, opts)
	if err != nil{
		log.Fatal("error creating new session ", err)
	}
}
//edited 17:38
func writeUser(id uuid.UUID)(reid interface{}) {
	log.Print(id)
	res, err := collection.InsertOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Fatalf("insertOne to the DB error: %s", err)
	}
	reid = res.InsertedID
	return reid
}
//edited 18:50
func readSessionInfo(id uuid.UUID) (info Session){
	filter := bson.D{{"_id", id}}
	err := collection.FindOne(context.Background(), filter, nil).Decode(&info)
	if err != nil {
		log.Fatalf("findOne to the DB error: %s", err)
	}
	return info
}
//TO EDIT 18:50
func delRefresh(RefreshT string, sess Session) {
//	how to read LATEST value
callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		insertNewSession(sess)
	_ = collection.FindOneAndDelete(context.TODO(), bson.M{"session": bson.M{"refresh":-1}}, nil)
	return nil, nil
}
session, err := client.StartSession() 
if err != nil{
	log.Fatalf("error starting session %s", err)
}
defer session.EndSession(context.Background())
_, _ = session.WithTransaction(context.Background(), callback)
//setTokens(ctx)
}

//TO EDIT 13/41
func delAllRefresh(id uuid.UUID) {

	_, err := collection.DeleteMany(context.Background(), bson.M{"_id": id}, nil)
	if err != nil {
		log.Fatalf("error deleting mane from the DB %s", err)
	}
}