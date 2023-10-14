module github.com/peterpasers/recipe-shop-bot

go 1.21.3

require github.com/peterparser/recipe-shop-bot/bot v1.0.0

require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 // indirect

replace github.com/peterparser/recipe-shop-bot/bot v1.0.0 => ./internal/bot
