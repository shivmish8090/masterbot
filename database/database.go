package database

import (
	"context"
	"log"
	"time"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	client  *mongo.Client
	userDB  *mongo.Collection
	chatDB  *mongo.Collection
)

func Init() {
	if config.MongoUri == "" {
		log.Panic("MongoDB URI is missing in config.MongoUri")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.MongoUri))
	if err != nil {
		log.Panicf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("EditGuardainBot")
	userDB = db.Collection("userstats")
	chatDB = db.Collection("chats")
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error while disconnecting MongoDB: %v", err)
	}
}

func IsServedUser(ctx context.Context, userID int) (bool, error) {
	count, err := userDB.CountDocuments(ctx, bson.M{"user_id": userID})
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

func IsServedChat(ctx context.Context, chatID int) (bool, error) {
	count, err := chatDB.CountDocuments(ctx, bson.M{"chat_id": chatID})
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