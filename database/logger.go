pacakge database



func IsLoggerEnabled() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var result struct {
		Enabled bool `bson:"enabled"`
	}

	err := loggerDB.FindOne(ctx, bson.M{}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return result.Enabled, nil
}


func SetLogger(enabled bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	update := bson.M{"$set": bson.M{"enabled": enabled}}
	opts := options.Update().SetUpsert(true)

	_, err := loggerDB.UpdateOne(ctx, bson.M{}, update, opts)
	return err
}
