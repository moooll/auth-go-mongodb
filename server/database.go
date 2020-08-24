package main

import (
	"log"
	"context"
		
	"github.com/gobuffalo/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"golang.org/x/crypto/bcrypt"
)
//User struct describes a user 
type User struct {
	ID uuid.UUID `bson: "_id`
	Session Session `bson: "session_id`
}

const dbURI =  "mongodb://localhost:27017"
//const dbURI = "mongodb+srv://goMongo:DBAuthApp@cluster0.p2liz.mongodb.net/authGo?retryWrites=true&w=majority&tlsInsecure=true"
var collection 	*mongo.Collection
//var cntx, _ = context.WithTimeout(context.TODO(), 10*time.Second)
//defer cancel()
var client mongo.Client

func connectToTheDB() {
	//context, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	client, err  := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI)) 
	//defer cancel()
	if err != nil {
		log.Fatal("problems connecting to the db:(",err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("authGo").Collection("users")
}

//check if user is in dB passing id as arg
func findUser(id uuid.UUID) (isPresent bool){
	//if user is present in DB isPresent is true
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}else {
		isPresent = true
	}
	return isPresent
}

func readUsers() (results User) {
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()){
		err = cur.Decode(results)
		if err != nil {
			log.Fatal(err)
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
		log.Fatal()
	}
}
//edited 17:38
func writeUser(id uuid.UUID) {
	_, err := collection.InsertOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}
}
//edited 18:50
func readSessionInfo(id uuid.UUID) (info Session){
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.Background(), filter, nil).Decode(&info)
	if err != nil {
		log.Fatal(err)
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
	log.Fatal(err)
}
defer session.EndSession(context.Background())
_, _ = session.WithTransaction(context.Background(), callback)
//setTokens(ctx)
}

//TO EDIT 13/41
func delAllRefresh(id uuid.UUID) {

	_, err := collection.DeleteMany(context.Background(), bson.M{"_id": id}, nil)
	if err != nil {
		log.Fatal(err)
	}
}