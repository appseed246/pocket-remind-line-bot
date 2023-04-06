package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"pocket-reminder/shared/convert"
	"pocket-reminder/shared/pocket"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"
	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("no IP in HTTP response")
	// ErrNon200Response non 200 status code in response
	ErrNon200Response        = errors.New("non 200 Response found")
	ErrNoConsumerKey         = errors.New("comsumer key is not set")
	ErrNoAccessToken         = errors.New("access token is not set")
	ErrNoChannelAccessToken  = errors.New("channel access token is not set")
	ErrNoChannelAccessSecret = errors.New("channel access secret is not set")
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
	fmt.Println("instance created.")

	// APIリクエストをパース
	lineEvents, err := parseAPIGatewayProxyRequest(&request)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	fmt.Println("request parsed.")

	// get pocket api keys
	consumerKey := os.Getenv("CONSUMER_KEY")
	accessToken := os.Getenv("ACCESS_TOKEN")

	for _, event := range lineEvents {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				command, parameter := parseMessage(message.Text)
				switch command {
				case "/carousel":
					fmt.Println("case: /carousel")
					response := pocket.FetchItems(consumerKey, accessToken)
					columns := convert.CreateCarouselMessage(response)
					template := linebot.NewTemplateMessage("Pocket Items", linebot.NewCarouselTemplate(columns...))
					if _, err = bot.ReplyMessage(event.ReplyToken, template).Do(); err != nil {
						fmt.Println(err)
					}
				case "/archive":
					fmt.Println("case: /archive")
					actions := []pocket.ModifyAction{{Action: "archive", ItemId: parameter[0]}}
					_, err := pocket.ModifyItem(consumerKey, accessToken, &actions)

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
					_, err := pocket.ModifyItem(consumerKey, accessToken, &actions)
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
