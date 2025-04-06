package modules 

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

