package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

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

func IsServedUser(userID int64) (bool, error) {
	v, ok := cache.Load("users")
	if ok {
		users := v.([]int64)
		for _, id := range users {
			if id == userID {
				return true, nil
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	count, err := userDB.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return false, err
	}
	if count > 0 {
		go func() {
			val, _ := cache.LoadOrStore("users", []int64{})
			users := val.([]int64)
			users = append(users, userID)
			cache.Store("users", users)
		}()
	}
	return count > 0, nil
}

func AddServedUser(userID int64) error {
	exists, err := IsServedUser(userID)
	if err != nil || exists {
		return err
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_, err := userDB.InsertOne(ctx, bson.M{"user_id": userID})
		if err == nil {
			val, _ := cache.LoadOrStore("users", []int64{})
			users := val.([]int64)
			users = append(users, userID)
			cache.Store("users", users)
		}
	}()
	return nil
}

func DeleteServedUser(userID int64) error {
	exists, err := IsServedUser(userID)
	if err != nil || !exists {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = userDB.DeleteOne(ctx, bson.M{"user_id": userID})
	if err == nil {
		val, _ := cache.Load("users")
		if val != nil {
			users := val.([]int64)
			for i, id := range users {
				if id == userID {
					users = append(users[:i], users[i+1:]...)
					break
				}
			}
			cache.Store("users", users)
		}
	}
	return err
}
