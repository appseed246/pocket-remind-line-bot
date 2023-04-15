package settings

import "errors"

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

	// TODO: 外部から取得
	RedirectURL     = "https://o9gzt4bf5k.execute-api.ap-northeast-1.amazonaws.com/Prod/token"
	AuthEndpointURL = "https://o9gzt4bf5k.execute-api.ap-northeast-1.amazonaws.com/Prod/auth"
)
