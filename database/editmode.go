package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

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
