package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	client   *mongo.Client
	userDB   *mongo.Collection
	chatDB   *mongo.Collection
	dbName   = "your_db_name"
	mongoURI = "your_mongo_uri"
	timeout  = 10 * time.Second
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	userDB = client.Database(dbName).Collection("userstats")
	chatDB = client.Database(dbName).Collection("chats")
	log.Println("Database initialized")
}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
	} else {
		log.Println("Database disconnected")
	}
}

// User Functions
func IsServedUser(ctx context.Context, userID int) (bool, error) {
	filter := bson.M{"user_id": userID}
	count, err := userDB.CountDocuments(ctx, filter)
	return count > 0, err
}

func GetServedUsers(ctx context.Context) ([]bson.M, error) {
	cursor, err := userDB.Find(ctx, bson.M{"user_id": bson.M{"$gt": 0}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []bson.M
	err = cursor.All(ctx, &users)
	return users, err
}

func AddServedUser(ctx context.Context, userID int) error {
	exists, err := IsServedUser(ctx, userID)
	if err != nil || exists {
		return err
	}
	_, err = userDB.InsertOne(ctx, bson.M{"user_id": userID})
	return err
}

func DeleteServedUser(ctx context.Context, userID int) error {
	_, err := userDB.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}

// Chat Functions
func IsServedChat(ctx context.Context, chatID int) (bool, error) {
	filter := bson.M{"chat_id": chatID}
	count, err := chatDB.CountDocuments(ctx, filter)
	return count > 0, err
}

func GetServedChats(ctx context.Context) ([]bson.M, error) {
	cursor, err := chatDB.Find(ctx, bson.M{"chat_id": bson.M{"$lt": 0}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []bson.M
	err = cursor.All(ctx, &chats)
	return chats, err
}

func AddServedChat(ctx context.Context, chatID int) error {
	exists, err := IsServedChat(ctx, chatID)
	if err != nil || exists {
		return err
	}
	_, err = chatDB.InsertOne(ctx, bson.M{"chat_id": chatID})
	return err
}

func DeleteServedChat(ctx context.Context, chatID int) error {
	_, err := chatDB.DeleteOne(ctx, bson.M{"chat_id": chatID})
	return err
}
