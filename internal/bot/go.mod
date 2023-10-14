module github.com/peterparser/recipe-shop-bot/bot

go 1.21.3

require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1

require github.com/peterparser/recipe-shop-bot/retriever v1.0.0

replace github.com/peterparser/recipe-shop-bot/retriever v1.0.0 => ../retriever
