package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogWhatsApp struct {
	ObjectId string  `json:"object_id" bson:"_id"`
	Error    string  `json:"error" bson:"error"`
	Method   string  `json:"method" bson:"method"`
	Endpoint string  `json:"endpoint" bson:"endpoint"`
	Payload  Payload `json:"payload" bson:"payload"`
}

type Payload struct {
	Email         string `json:"email" bson:"email"`
	Password      string `json:"password" bson:"password"`
	Telepon       string `json:"telepon" bson:"telepon"`
	Nama          string `json:"nama" bson:"nama"`
	Nama_Lengkap  string `json:"nama_lengkap" bson:"nama_lengkap"`
	Nama_Depan    string `json:"nama_depan" bson:"nama_depan"`
	Nama_Belakang string `json:"nama_belakang" bson:"nama_belakang"`
	Kategori      string `json:"kategori" bson:"kategori"`
	Gender        string `json:"gender" bson:"gender"`
	Device        string `json:"device" bson:"device"`
	Otp           string `json:"otp" bson:"otp"`
}

func main() {
	// Set up the MongoDB client and connect to the database
	// mongodb://ocistok:2BelasJuta$@dds-d9j719c71af16d941168-pub.mongodb.ap-southeast-5.rds.aliyuncs.com:3717,dds-d9j719c71af16d942147-pub.mongodb.ap-southeast-5.rds.aliyuncs.com:3717,dds-d9j719c71af16d943156-pub.mongodb.ap-southeast-5.rds.aliyuncs.com:3717/ocistok?replicaSet=mgset-1104029287&authSource=admin

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb://ocistok:2BelasJuta$@dds-d9j719c71af16d941168-pub.mongodb.ap-southeast-5.rds.aliyuncs.com:3717,dds-d9j719c71af16d942147-pub.mongodb.ap-southeast-5.rds.aliyuncs.com:3717,dds-d9j719c71af16d943156-pub.mongodb.ap-southeast-5.rds.aliyuncs.com:3717/ocistok?replicaSet=mgset-1104029287&authSource=admin").SetServerAPIOptions(serverAPI)

	//create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	// Select the MongoDB database and collection
	collection := client.Database("whatsappScrap").Collection("logwhatsapp")

	// Define a filter to find documents
	filter := bson.M{}

	// Retrieve documents from the collection
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	// Iterate through the documents and extract the timestamp from the _id field
	for cur.Next(context.Background()) {
		var doc LogWhatsApp
		if err := cur.Decode(&doc); err != nil {
			log.Fatal(err)
		}

		objectID, err := primitive.ObjectIDFromHex(doc.ObjectId)
		if err != nil {
			log.Fatal(err)
		}
		timestamp := objectID.Timestamp()

		date := time.Unix(int64(timestamp.Unix()), 0)

		//update
		extractedId, err := primitive.ObjectIDFromHex(doc.ObjectId)
		fmt.Println(doc.ObjectId, date, extractedId)
		if err != nil {
			log.Fatal(err)
		}
		filterforupdate := bson.D{{"_id", extractedId}}
		update := bson.D{{"$set", bson.D{{"date", date}}}}
		collection.UpdateOne(context.TODO(), filterforupdate, update)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Disconnect from MongoDB
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}
