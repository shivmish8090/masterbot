package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EchoSettings struct {
	ChatID int64  `bson:"chat_id"`
	Mode   string `bson:"mode"`
	Limit  int    `bson:"limit"`
}

const (
	defaultLimit = 800
	defaultMode  = "MANUAL"
)

func SetEchoSettings(data *EchoSettings) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"mode":  data.Mode,
			"limit": data.Limit,
		},
	}

	_, err := echoDB.UpdateOne(ctx, bson.M{"chat_id": data.ChatID}, update, options.UpdateOne().SetUpsert(true))
	return err
}

func GetEchoSettings(chatID int64) (*EchoSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var settings EchoSettings
	err := echoDB.FindOne(ctx, bson.M{"chat_id": chatID}).Decode(&settings)
	if err != nil {
		return &EchoSettings{
			ChatID: chatID,
			Mode:   defaultMode,
			Limit:  defaultLimit,
		}, nil
	}

	if settings.Mode == "" {
		settings.Mode = defaultMode
	}
	if settings.Limit == 0 {
		settings.Limit = defaultLimit
	}

	return &settings, nil
}
