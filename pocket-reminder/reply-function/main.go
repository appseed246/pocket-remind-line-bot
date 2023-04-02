package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"pocket-reminder/shared/convert"
	"pocket-reminder/shared/pocket"

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

	// カルーセル作成
	response := pocket.FetchItems(os.Getenv("CONSUMER_KEY"), os.Getenv("ACCESS_TOKEN"))
	columns := convert.CreateCarouselMessage(response)
	template := linebot.NewTemplateMessage("Pocket Items", linebot.NewCarouselTemplate(columns...))

	for _, event := range lineEvents {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "carousel" {
					if _, err = bot.ReplyMessage(event.ReplyToken, template).Do(); err != nil {
						log.Print(err)
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
