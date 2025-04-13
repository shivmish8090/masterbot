package helpers


func GetAdmins(b Bot, ChatId int64) ([]int64, error) {
	cacheKey := fmt.Sprintf("admins:%d", ChatId)

	if admins, ok := config.LoadTyped[[]int64](config.Cache, cacheKey); ok {
		return admins, nil
	}

	chatmembers, err := b.GetChatAdministrators(ChatId, nil)
	if err != nil {
		return nil, err
	}

	var admins []int64
	for _, m := range chatmembers {
		status := m.GetStatus()
		if status == "administrator" || status == "creator" {
			admins = append(admins, m.GetUser().Id)
		}
	}

	config.Cache.Store(cacheKey, admins)

	return admins, nil
}