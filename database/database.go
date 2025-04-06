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

func IsServedUser(userID int64) (bool, error) {
	key := fmt.Sprintf("users:%d", userID)
	if _, ok := cache.Load(key); ok {
		return true, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	count, err := userDB.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return false, err
	}

	if count > 0 {
		cache.Store(key, true)
	}

	return count > 0, nil
}

func GetServedUsers() ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cursor, err := userDB.Find(ctx, bson.M{"user_id": bson.M{"$gt": 0}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []struct {
		UserID int64 `bson:"user_id"`
	}

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	var userIDs []int64
	for _, u := range users {
		userIDs = append(userIDs, u.UserID)
	}

	return userIDs, nil
}

func AddServedUser(userID int64) error {
	key := fmt.Sprintf("users:%d", userID)
	if _, ok := cache.Load(key); !ok {
		exists, err := IsServedUser(userID)
		if err != nil || exists {
			return err
		}
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_, err := userDB.InsertOne(ctx, bson.M{"user_id": userID})
		if err == nil {
			cache.Store(key, true)
		}
	}()
	return nil
}

func DeleteServedUser(userID int64) error {
	key := fmt.Sprintf("users:%d", userID)
	if _, ok := cache.Load(key); !ok {
		exists, err := IsServedUser(userID)
		if err != nil || !exists {
			return err
		}
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_, err := userDB.DeleteOne(ctx, bson.M{"user_id": userID})
		if err == nil {
			cache.Delete(key)
		}
	}()
	return nil
}

func IsServedChat(chatID int64) (bool, error) {
	key := fmt.Sprintf("chats:%d", chatID)
	if _, ok := cache.Load(key); ok {
		return true, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	count, err := chatDB.CountDocuments(ctx, bson.M{"chat_id": chatID})
	if err != nil {
		return false, err
	}

	if count > 0 {
		cache.Store(key, true)
	}

	return count > 0, nil
}

func GetServedChats() ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cursor, err := chatDB.Find(ctx, bson.M{"chat_id": bson.M{"$lt": 0}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []struct {
		ChatID int64 `bson:"chat_id"`
	}
	if err = cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	var chatIDs []int64
	for _, chat := range chats {
		chatIDs = append(chatIDs, chat.ChatID)
	}

	return chatIDs, nil
}

func AddServedChat(chatID int64) error {
	key := fmt.Sprintf("chats:%d", chatID)
	if _, ok := cache.Load(key); !ok {
		exists, err := IsServedChat(chatID)
		if err != nil || exists {
			return err
		}
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_, err := chatDB.InsertOne(ctx, bson.M{"chat_id": chatID})
		if err == nil {
			cache.Store(key, true)
		}
	}()
	return nil
}

func DeleteServedChat(chatID int64) error {
	key := fmt.Sprintf("chats:%d", chatID)
	if _, ok := cache.Load(key); !ok {
		exists, err := IsServedChat(chatID)
		if err != nil || !exists {
			return err
		}
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_, err := chatDB.DeleteOne(ctx, bson.M{"chat_id": chatID})
		if err == nil {
			cache.Delete(key)
		}
	}()
	return nil
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
