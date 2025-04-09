package database

import (
	"context"
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

	// Unique index on userstats.user_id
	_, err = userDB.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"user_id": 1},
		Options: options.Index().SetName("unique_user_id").SetUnique(true),
	})
	if err != nil {
		log.Printf("Failed to create unique index on userstats: %v", err)
	}

	// Unique index on chats.chat_id
	_, err = chatDB.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"chat_id": 1},
		Options: options.Index().SetName("unique_chat_id").SetUnique(true),
	})
	if err != nil {
		log.Printf("Failed to create unique index on chats: %v", err)
	}

	// Unique index on editmodes.chat_id
	_, err = editModeDB.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"chat_id": 1},
		Options: options.Index().SetName("unique_chat_id_editmode").SetUnique(true),
	})
	if err != nil {
		log.Printf("Failed to create unique index on editmodes: %v", err)
	}

	// Unique index on echos.chat_id
	_, err = echoDB.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"chat_id": 1},
		Options: options.Index().SetName("unique_chat_id_echo").SetUnique(true),
	})
	if err != nil {
		log.Printf("Failed to create unique index on echos: %v", err)
	}
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Error while disconnecting MongoDB: %v", err)
	}
}
