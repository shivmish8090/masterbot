package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

var (
	client     *mongo.Client
	userDB     *mongo.Collection
	chatDB     *mongo.Collection
	editModeDB *mongo.Collection
	echoDB     *mongo.Collection
	timeout    = 10 * time.Second
)

func init() {
	if config.MongoUri == "" {
		log.Panic("MongoDB URI is missing in config.MongoUri")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var err error
	client, err = mongo.Connect(options.Client().ApplyURI(config.MongoUri))
	if err != nil {
		log.Panicf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("EditGuardainBot")
	userDB = db.Collection("userstats")
	chatDB = db.Collection("chats")
	editModeDB = db.Collection("editmodes")
	echoDB = db.Collection("echos")

	_, err = userDB.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"user_id": 1}})
	if err != nil {
		log.Printf("Failed to create index on userstats: %v", err)
	}

	_, err = chatDB.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"chat_id": 1}})
	if err != nil {
		log.Printf("Failed to create index on chats: %v", err)
	}
	_, err = editModeDB.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"chat_id": 1}})
	if err != nil {
		log.Printf("Failed to create index on editmodes: %v", err)
	}
	_, err = echoDB.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"chat_id": 1},
	})
	if err != nil {
		log.Printf("Failed to create index on echos: %v", err)
	}
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error while disconnecting MongoDB: %v", err)
	}
}


// SetEditMode sets the edit mode for a chat ("ADMIN", "USER", "OFF").
func SetEditMode(chatID int64, mode string) (bool, error) {
	key := fmt.Sprintf("editmode:%d", chatID)

	if cachedMode, ok := cache.Load(key); ok {
		if strMode, valid := cachedMode.(string); valid && strMode == mode {
			return true, nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := editModeDB.UpdateOne(
		ctx,
		bson.M{"chat_id": chatID},
		bson.M{"$set": bson.M{"mode": mode}},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		log.Printf("SetEditMode error for chatID %d: %v", chatID, err)
		return false, err
	}

	if result.ModifiedCount > 0 || result.UpsertedCount > 0 {
		cache.Store(key, mode)
		return true, nil
	}

	cache.Store(key, mode)
	return false, nil
}

func GetEditMode(chatID int64) string {
	key := fmt.Sprintf("editmode:%d", chatID)

	if value, ok := cache.Load(key); ok {
		if mode, valid := value.(string); valid {
			return mode
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var result struct {
		Mode string `bson:"mode"`
	}

	err := editModeDB.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&result)
	if err != nil || result.Mode == "" {
		return "USER"
	}

	cache.Store(key, result.Mode)
	return result.Mode
}
