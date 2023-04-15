package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"pocket-reminder/shared/convert"
	"pocket-reminder/shared/datasource"
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
		log.Fatal("failed to create line client.", err)
	}

	// DbClient初期化
	db, err := datasource.New()
	if err != nil {
		log.Fatal("datasource create error.", err)
	}

	// TODO: 現在は自分専用
	lineUserIdList := []string{os.Getenv("OWN_LINE_ID")}

	// TODO: Botを利用するすべてのユーザへのブロードキャンストに対応
	for _, lineUserId := range lineUserIdList {
		user, err := db.GetPocketReminderUser(lineUserId)
		if err != nil {
			log.Println("user info access error: " + err.Error())
			break
		}
		if user == nil {
			log.Println("user not found.")
			break
		}

		pocketClient, err := pocket.New(
			os.Getenv("CONSUMER_KEY"),
			user.PocketAccessToken,
		)
		if err != nil {
			log.Fatal("failed to create pocket client.", err)
		}

		// カルーセル作成
		response := pocketClient.FetchItems()
		columns := convert.CreateCarouselMessage(response)
		convert.PrintCarouselColumns(columns)
		template := linebot.NewTemplateMessage("Pocket Items", linebot.NewCarouselTemplate(columns...))

		// メッセージをブロードキャスト
		if _, err := bot.PushMessage(lineUserId, template).Do(); err != nil {
			log.Printf("Failed to send message : %v", err)
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
