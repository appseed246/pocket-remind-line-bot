package main

import (
	"errors"
	"log"
	"net/url"
	"os"
	"pocket-reminder/shared/pocket"
	"pocket-reminder/shared/settings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lineUserId := request.QueryStringParameters["lineUserId"]
	if lineUserId == "" {
		log.Println("lineUserId not found")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, errors.New("lineUserId not found")
	}

	// Pocket: oauth/request
	authorizationCode, err := pocket.GetAuthorizationCode(os.Getenv("CONSUMER_KEY"), settings.RedirectURL)
	if err != nil {
		log.Println("get authorization code failed")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// redirect to Pocket Login
	params := url.Values{}
	params.Set("request_token", authorizationCode)
	params.Set("redirect_uri", settings.RedirectURL+"?code="+authorizationCode+"&lineUserId="+lineUserId)
	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Location": pocket.BASE_URL + pocket.PATH_OAUTH_AUTHZ + "?" + params.Encode(),
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
