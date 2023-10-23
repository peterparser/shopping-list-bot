module github.com/peterparser/recipe-shop-bot/bot

go 1.21.3

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/peterparser/recipe-shop-bot/retriever v0.0.0-00010101000000-000000000000
)

replace github.com/peterparser/recipe-shop-bot/retriever => ../retriever
