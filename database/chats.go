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
