package telegramapi

import (
	"fmt"
	"strings"

	"github.com/Samar2170/portfolio-manager/account"
)

func TgRegisterDematAccount(lines []string, chatID int64) error {
	payload := make(map[string]string)
	for _, l := range lines {
		words := strings.Split(l, ":")
		if len(words) != 2 {
			return fmt.Errorf("each line should contain param and value")
		}
		payload[words[0]] = words[1]
	}
	if _, ok := payload["code"]; !ok {
		return fmt.Errorf("payload code doesnt exist")
	}
	if _, ok := payload["broker"]; !ok {
		return fmt.Errorf("broker doesnt exist")
	}

	user, err := account.GetUserByChatId(chatID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	da := account.DematAccount{UserId: user.ID, Code: payload["code"], Broker: payload["broker"]}
	da.Create()
	return nil
}
