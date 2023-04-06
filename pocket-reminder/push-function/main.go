package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"pocket-reminder/shared/convert"
	"pocket-reminder/shared/pocket"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	fmt.Println(event.Detail)

	// LineClient初期化
	bot, err := linebot.New(
		os.Getenv("CHANNEL_ACCESS_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal("failed to create line client.")
	}

	// Pocket API Client初期化
	pocketClient, err := pocket.New(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal("failed to create pocket client.")
	}

	fmt.Println("instance created.")

	// カルーセル作成
	response := pocketClient.FetchItems()
	columns := convert.CreateCarouselMessage(response)
	convert.PrintCarouselColumns(columns)
	template := linebot.NewTemplateMessage("Pocket Items", linebot.NewCarouselTemplate(columns...))

	// メッセージをブロードキャスト
	if _, err := bot.BroadcastMessage(template).Do(); err != nil {
		log.Printf("Failed to send message : %v", err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
