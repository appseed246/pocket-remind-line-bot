package main

import (
	"encoding/json"
	"fmt"
	"os"
	"pocket-reminder/shared/convert"
	"pocket-reminder/shared/datasource"
	"pocket-reminder/shared/pocket"
	"pocket-reminder/shared/settings"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("**** body ****")
	fmt.Println(request.Body)
	fmt.Println("**** body ****")

	// LineClient初期化
	bot, err := linebot.New(
		os.Getenv("CHANNEL_ACCESS_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	db, err := datasource.New()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	fmt.Println("instance created.")

	// APIリクエストをパース
	lineEvents, err := parseAPIGatewayProxyRequest(&request)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	fmt.Println("request parsed.")

	for _, event := range lineEvents {
		if event.Type == linebot.EventTypeMessage {
			lineUserId := event.Source.UserID

			// PocketAPIのアクセストークン取得
			// ユーザ情報が存在しない場合ログインを要求する
			user, err := db.GetPocketReminderUser(lineUserId)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			if user == nil || user.PocketAccessToken == "" {
				if _, err = bot.ReplyMessage(event.ReplyToken, createLoginRequestMessage(lineUserId)).Do(); err != nil {
					fmt.Println(err)
				}
				break
			}
			accessToken := user.PocketAccessToken

			// Pocket API Client初期化
			pocketClient, err := pocket.New(
				os.Getenv("CONSUMER_KEY"),
				accessToken,
			)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				command, parameter := parseMessage(message.Text)
				switch command {
				case "/auth":
					fmt.Println("case: /auth")
					if _, err = bot.ReplyMessage(event.ReplyToken, createLoginRequestMessage(lineUserId)).Do(); err != nil {
						fmt.Println(err)
					}
				case "/carousel":
					fmt.Println("case: /carousel")
					response := pocketClient.FetchItems()
					columns := convert.CreateCarouselMessage(response)
					template := linebot.NewTemplateMessage("Pocket Items", linebot.NewCarouselTemplate(columns...))
					if _, err = bot.ReplyMessage(event.ReplyToken, template).Do(); err != nil {
						fmt.Println(err)
					}
				case "/archive":
					fmt.Println("case: /archive")
					actions := []pocket.ModifyAction{{Action: "archive", ItemId: parameter[0]}}
					_, err := pocketClient.ModifyItem(&actions)

					if err != nil {
						fmt.Println(err.Error())
					} else {
						confirm := linebot.NewConfirmTemplate(
							"アイテムをアーカイブしました。",
							linebot.NewMessageAction("元に戻す", fmt.Sprintf("/readd %s", parameter[0])),
							linebot.NewMessageAction("リストを見る", "/carousel"),
						)
						template := linebot.NewTemplateMessage("アイテムをアーカイブしました。", confirm)
						if _, err = bot.ReplyMessage(event.ReplyToken, template).Do(); err != nil {
							fmt.Println(err)
						}

					}
				case "/readd":
					fmt.Println("case: /readd")
					actions := []pocket.ModifyAction{{Action: "readd", ItemId: parameter[0]}}
					_, err := pocketClient.ModifyItem(&actions)
					if err != nil {
						fmt.Println(err.Error())
					} else {
						confirm := linebot.NewConfirmTemplate(
							"アーカイブをキャンセルしました。",
							linebot.NewMessageAction("アーカイブ", fmt.Sprintf("/archive %s", parameter[0])),
							linebot.NewMessageAction("リストを見る", "/carousel"),
						)
						template := linebot.NewTemplateMessage("アーカイブをキャンセルしました。", confirm)
						if _, err = bot.ReplyMessage(event.ReplyToken, template).Do(); err != nil {
							fmt.Println(err)
						}

					}
				}
			}
		}
	}
	fmt.Println("replied.")

	return events.APIGatewayProxyResponse{
		Body:       "Hello, world",
		StatusCode: 200,
	}, nil
}

func createLoginRequestMessage(lineUserId string) *linebot.TemplateMessage {
	button := linebot.NewButtonsTemplate(
		"",
		"",
		"このアプリを利用するにはPocketにログインしてください。",
		linebot.NewURIAction("ログイン", settings.AuthEndpointURL+"?lineUserId="+lineUserId),
	)
	return linebot.NewTemplateMessage("ログイン依頼", button)
}

func parseMessage(message string) (string, []string) {
	parsed_message := strings.Split(message, " ")
	command := parsed_message[0]
	parameters := parsed_message[1:]

	return command, parameters
}

func parseAPIGatewayProxyRequest(r *events.APIGatewayProxyRequest) ([]*linebot.Event, error) {
	body := []byte(r.Body)

	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}

	if err := json.Unmarshal(body, request); err != nil {
		return nil, err
	}

	return request.Events, nil
}

func main() {
	lambda.Start(handler)
}
