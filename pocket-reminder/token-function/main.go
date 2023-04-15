package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"pocket-reminder/shared/datasource"
	"pocket-reminder/shared/pocket"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	println(request.QueryStringParameters)
	lineUserId := request.QueryStringParameters["lineUserId"]
	authorizationCode := request.QueryStringParameters["code"]
	if authorizationCode == "" {
		log.Println("authorization code not found")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, errors.New("authorization code not found")
	}
	if lineUserId == "" {
		log.Println("lineUserId not found")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, errors.New("lineUserId not found")
	}

	accessToken, userName, err := pocket.GetAccessToken(os.Getenv("CONSUMER_KEY"), authorizationCode)
	if err != nil {
		log.Println("authorization code not found")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, errors.New("access token get failed")
	}
	fmt.Println("get access token success. username = " + userName)

	// PocketAPIの情報を保存
	db, err := datasource.New()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// ユーザIDの生成 or 取得
	user, err := db.GetPocketReminderUser(lineUserId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	userId := ""
	if user == nil {
		uuidV4, err := uuid.NewRandom()
		if err != nil {
			log.Println("Error generating UUID:", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, err
		}
		userId = uuidV4.String()
	} else {
		userId = user.UserId
	}

	err = db.SaveAccessToken(userId, lineUserId, accessToken)
	if err != nil {
		log.Println("failed to save token")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, errors.New("failed to save token")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Location": "https://line.me/R/",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
