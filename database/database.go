package database

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

var (
	client  *mongo.Client
	userDB  *mongo.Collection
	chatDB  *mongo.Collection
	cache   sync.Map
	timeout = 10 * time.Second
)

func init() {
	if config.MongoUri == "" {
		log.Panic("MongoDB URI is missing in config.MongoUri")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.MongoUri))
	if err != nil {
		log.Panicf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("EditGuardainBot")
	userDB = db.Collection("userstats")
	chatDB = db.Collection("chats")

	_, err = userDB.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"user_id": 1}})
	if err != nil {
		log.Printf("Failed to create index on userstats: %v", err)
	}

	_, err = chatDB.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"chat_id": 1}})
	if err != nil {
		log.Printf("Failed to create index on chats: %v", err)
	}
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error while disconnecting MongoDB: %v", err)
	}
}

func IsServedUser(ctx context.Context, userID int) (bool, error) {
	if _, ok := cache.Load(userID); ok {
		return true, nil
	}

	count, err := userDB.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return false, err
	}
	if count > 0 {
		cache.Store(userID, true)
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
	if err != nil || exists {
		return err
	}

	go func() {
		_, err := userDB.InsertOne(context.Background(), bson.M{"user_id": userID})
		if err == nil {
			cache.Store(userID, true)
		}
	}()
	return nil
}

func DeleteServedUser(ctx context.Context, userID int) error {
	_, err := userDB.DeleteOne(ctx, bson.M{"user_id": userID})
	if err == nil {
		cache.Delete(userID)
	}
	return err
}

func IsServedChat(ctx context.Context, chatID int) (bool, error) {
	if _, ok := cache.Load(chatID); ok {
		return true, nil
	}

	count, err := chatDB.CountDocuments(ctx, bson.M{"chat_id": chatID})
	if err != nil {
		return false, err
	}
	if count > 0 {
		cache.Store(chatID, true)
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
	if err != nil || exists {
		return err
	}

	go func() {
		_, err := chatDB.InsertOne(context.Background(), bson.M{"chat_id": chatID})
		if err == nil {
			cache.Store(chatID, true)
		}
	}()
	return nil
}

func DeleteServedChat(ctx context.Context, chatID int) error {
	_, err := chatDB.DeleteOne(ctx, bson.M{"chat_id": chatID})
	if err == nil {
		cache.Delete(chatID)
	}
	return err
}
