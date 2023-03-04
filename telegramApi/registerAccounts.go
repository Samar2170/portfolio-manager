package telegramapi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Samar2170/portfolio-manager/account"
)

func TgRegisterUser(msg string, chatId int64) string {
	words := strings.Split(msg, " ")
	otp, err := strconv.ParseInt(words[1], 10, 64)
	if err != nil {
		return "otp is not an integer"
	}
	user, err := account.GetUserByTelegramOTP(uint(otp))
	if err != nil {
		return "user not found for otp"
	}
	tgUser := account.TelegramUser{User: user, ChatId: chatId}
	tgUser.Create()
	return "user created successfully"
}

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
func TgRegisterBankAccount(lines []string, chatID int64) error {
	payload := make(map[string]string)
	for _, l := range lines {
		words := strings.Split(l, ":")
		if len(words) != 2 {
			return fmt.Errorf("each line should contain param and value")
		}
		payload[words[0]] = words[1]
	}
	if _, ok := payload["account_number"]; !ok {
		return fmt.Errorf("payload account_number doesnt exist")
	}
	if _, ok := payload["bank"]; !ok {
		return fmt.Errorf("bank doesnt exist")
	}

	user, err := account.GetUserByChatId(chatID)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	ba := account.BankAccount{UserId: user.ID, AccountNo: payload["account_number"], Bank: payload["bank"]}
	ba.Create()
	return nil
}
