package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/telegramApi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

var Tgbot *tgbotapi.BotAPI
var BotToken string
var APIHost string
var Port string

type APIClient struct {
	Token string
}

func (ac APIClient) getToken(chatID int64) error {
	time.Sleep(5 * time.Second)
	chatIDStr := strconv.FormatInt(chatID, 10)
	resp, err := http.PostForm(fmt.Sprintf("http://%s:%s/%s", APIHost, Port, "chat-id-login"), url.Values{"chat_id": {chatIDStr}})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp)
	return nil
}

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	BotToken = viper.GetString("BOTTOKEN")
	APIHost = viper.GetString("HOST")
	Port = viper.GetString("PORT")
}

func StartBot() {
	var err error
	Tgbot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Println(err.Error())
	}
	Tgbot.Debug = true
	log.Printf("Authorized on account %s", Tgbot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := Tgbot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		textMsg := MsgRouter(update)
		msg := tgbotapi.NewMessage(chatID, textMsg)
		Tgbot.Send(msg)
	}
}

func parseForm(msg string, chatID int64) string {
	lines := strings.Split(msg, "\n")
	for i, l := range lines {
		if string(l[0]) == "[" && string(l[len(l)-1]) == "]" {
			formName := string(l[1 : len(l)-1])
			switch formName {
			case "RegisterBankAccounts":
				telegramApi.TgRegisterDematAccount(lines[i:], chatID)
			case "RegisterDematAccounts":
				telegramApi.TgRegisterDematAccount(lines[i:], chatID)
			}
		} else {
			continue
		}

	}
	return strings.Join(lines, ",")
}
func registerUser(msg string, chatId int64) string {
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

func MsgRouter(update tgbotapi.Update) string {
	msg := update.Message.Text
	lines := strings.Split(msg, "\n")
	route := strings.Split(lines[0], " ")
	fmt.Println(route[0])
	switch route[0] {
	case "/register":
		return registerUser(update.Message.Text, update.Message.Chat.ID)
	case "/register-trade":
		return parseForm(update.Message.Text, update.Message.Chat.ID)
	case "/subscribe":
		if len(route[0]) < 2 {
			return "Please specify a service"
		}
		return fmt.Sprintf(update.Message.Chat.FirstName, update.Message.Chat.LastName, update.Message.Chat.ID)
	case "/help":
		return "Help"
	case "/todo":
		return "Todo"
	default:
		return "Unknown command"
	}
}
