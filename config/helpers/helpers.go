package helpers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

type AdminData struct {
	Status      string
	Permissions gotgbot.MergedChatMember
}

func FetchAdmins(b gotgbot.Bot, ChatId int64) (map[int64]AdminData, error) {
	cacheKey := fmt.Sprintf("admins:%d", ChatId)

	if admins, ok := config.LoadTyped[map[int64]AdminData](config.Cache, cacheKey); ok {
		return admins, nil
	}

	chatmembers, err := b.GetChatAdministrators(ChatId, nil)
	if err != nil {
		return nil, err
	}

	adminMap := make(map[int64]AdminData)
	for _, m := range chatmembers {
		status := m.GetStatus()
		if status == "administrator" || status == "creator" {
			adminMap[m.GetUser().Id] = AdminData{
				Status:      status,
				Permissions: m.MergeChatMember(),
			}
		}
	}

	config.Cache.Store(cacheKey, adminMap)
	return adminMap, nil
}

func GetAdmins(b gotgbot.Bot, ChatId int64) ([]int64, error) {
	adminMap, err := FetchAdmins(b, ChatId)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for id := range adminMap {
		ids = append(ids, id)
	}

	return ids, nil
}

func GetOwner(b gotgbot.Bot, ChatId int64) (int64, error) {
	adminMap, err := FetchAdmins(b, ChatId)
	if err != nil {
		return 0, err
	}

	for id, data := range adminMap {
		if data.Status == "creator" {
			return id, nil
		}
	}

	return 0, fmt.Errorf("no creator found")
}
