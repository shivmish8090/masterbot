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

func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	userDB = client.Database(dbName).Collection("userstats")
	chatDB = client.Database(dbName).Collection("chats")
}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
	}
}

// User Functions
func IsServedUser(ctx context.Context, userID int) (bool, error) {
	filter := bson.M{"user_id": userID}
	count, err := userDB.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetServedUsers(ctx context.Context) ([]bson.M, error) {
	cursor, err := userDB.Find(ctx, bson.M{"user_id": bson.M{"$gt": 0}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []bson.M
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func AddServedUser(ctx context.Context, userID int) error {
	exists, err := IsServedUser(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return nil
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
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetServedChats(ctx context.Context) ([]bson.M, error) {
	cursor, err := chatDB.Find(ctx, bson.M{"chat_id": bson.M{"$lt": 0}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []bson.M
	if err = cursor.All(ctx, &chats); err != nil {
		return nil, err
	}
	return chats, nil
}

func AddServedChat(ctx context.Context, chatID int) error {
	exists, err := IsServedChat(ctx, chatID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = chatDB.InsertOne(ctx, bson.M{"chat_id": chatID})
	return err
}

func DeleteServedChat(ctx context.Context, chatID int) error {
	_, err := chatDB.DeleteOne(ctx, bson.M{"chat_id": chatID})
	return err
}
