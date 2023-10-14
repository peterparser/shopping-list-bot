package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/peterparser/recipe-shop-bot/retriever"
)

type User struct {
	Id   int64
	Link string
	Gid  int
}

var users map[int64]*User

func Start() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	users = make(map[int64]*User)

	if err != nil {
		panic(err)
	}
	log.Println("Bot started...")

	handlers := map[string]func(*tgbotapi.BotAPI, tgbotapi.Update){
		"config": handleConfig,
		"help":   handleHelp,
		"list":   handleList,
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}
		user, ok := users[update.Message.Chat.ID]

		if ok {
			if user.Link == "" {
				handleLink(bot, update, user)
				continue
			}
			if user.Gid == 0 {
				handleGid(bot, update, user)
				log.Printf("User %d configured!\n", user.Id)
				continue
			}

		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		handlers[update.Message.Command()](bot, update)
	}
}

func handleConfig(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send me the link of your google sheet")
	users[update.Message.Chat.ID] = &User{Id: update.Message.Chat.ID}
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func handleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I understand only /config and /list")
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func handleLink(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *User) {
	link := update.Message.Text
	user.Link = link
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ok, now send me the GID of the second workbook")
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}

}

func handleGid(bot *tgbotapi.BotAPI, update tgbotapi.Update, user *User) {
	gid := update.Message.Text
	gidNumber, err := strconv.Atoi(gid)
	if err != nil {
		log.Fatal(err)
	}
	user.Gid = gidNumber
}

func handleList(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	user := users[update.Message.Chat.ID]
	planning := retriever.RetrieveDoc(user.Link, 0)
	recipes := retriever.RetrieveDoc(user.Link, user.Gid)
	parsedPlan := retriever.ExtractDishFromPlan(planning)
	parsedRecipes := retriever.BuildRecipes(recipes)

	toBuy := make([]string, 0)

	for _, itemInPlan := range parsedPlan {
		ingredients, ok := parsedRecipes[itemInPlan]
		if ok {
			toBuy = append(toBuy, ingredients...)
		}
	}
	counter := make(map[string]int)
	for _, element := range toBuy {
		counter[element] += 1
	}
	messagesList := make([]string, 0)

	for item, quantity := range counter {
		messagesList = append(messagesList, fmt.Sprintf("%s x%d", item, quantity))
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(messagesList, "\n"))
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
