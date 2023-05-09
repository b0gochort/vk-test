package bot

import (
	"math/rand"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
)

func SendMessage(vk *api.VK, userID int, message string) error {
	// Создаем параметры запроса
	sendParams := params.NewMessagesSendBuilder()
	rand.Seed(time.Now().UnixNano())
	sendParams.RandomID(rand.Int())
	sendParams.PeerID(userID)
	sendParams.Message(message)

	// Отправляем сообщение
	_, err := vk.MessagesSend(sendParams.Params)
	if err != nil {
		return err
	}

	return nil
}

// func getLastBotMessage(vk *api.VK, peerID int) (*object.Message, error) {
// 	history, err := vk.MessagesGetHistory(api.Params{
// 		"peer_id":        peerID,
// 		"count":          1,
// 		"rev":            1,
// 		"is_out":         1,
// 		"extended":       1,
// 		"fields":         "id",
// 		"group_messages": 1,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(history.Items) == 0 {
// 		return nil, errors.New("no messages found")
// 	}
// 	for _, msg := range history.Items {
// 		if msg.FromID == vk.GroupID && msg.ID == history.Profiles[0].ID {
// 			return &msg, nil
// 		}
// 	}
// 	return nil, errors.New("bot message not found")
// }

func ConcatenateStrings(str []string) string {
	return strings.Join(str, "\n")
}

func GetTimeFromString(dateTimeStr string) (string, error) {
	dateTime, err := time.Parse(time.RFC3339, dateTimeStr)
	if err != nil {
		return "", err
	}
	timeStr := dateTime.Format("15:04")
	return timeStr, nil
}
