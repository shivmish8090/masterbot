package database

import (
        "context"
        "time"

        "go.mongodb.org/mongo-driver/v2/bson"
        "go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EchoSettings struct {
        ChatID int64  `bson:"chat_id"`
        Mode   string `bson:"mode"`   // OFF, MANUAL, or AUTOMATIC
        Limit  int    `bson:"limit"`  // Default 800
}

const (
        defaultLimit = 800
        defaultMode  = "MANUAL"
)

// SetEchoMode sets the echo mode (OFF, MANUAL, AUTOMATIC) for a chat
func SetEchoMode(chatID int64, mode string) error {
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()

        update := bson.M{
                "$set": bson.M{
                        "mode": mode,
                },
        }
        _, err := echoDB.UpdateOne(ctx, bson.M{"chat_id": chatID}, update, options.Update().SetUpsert(true))
        return err
}

// SetEchoLimit sets the character limit for echo handling
func SetEchoLimit(chatID int64, limit int) error {
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()

        update := bson.M{
                "$set": bson.M{
                        "limit": limit,
                },
        }
        _, err := echoDB.UpdateOne(ctx, bson.M{"chat_id": chatID}, update, options.Update().SetUpsert(true))
        return err
}

// GetEchoSettings returns the echo settings for a given chat.
// If not found, returns default values (MANUAL mode, 800 limit).
func GetEchoSettings(chatID int64) (*EchoSettings, error) {
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()

        var settings EchoSettings
        err := echoDB.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&settings)
        if err != nil {
                if err.Error() == "mongo: no documents in result" {
                        return &EchoSettings{
                                ChatID: chatID,
                                Mode:   defaultMode,
                                Limit:  defaultLimit,
                        }, nil
                }
                return nil, err
        }

        // Fallback to default if fields are missing
        if settings.Mode == "" {
                settings.Mode = defaultMode
        }
        if settings.Limit == 0 {
                settings.Limit = defaultLimit
        }

        return &settings, nil
}