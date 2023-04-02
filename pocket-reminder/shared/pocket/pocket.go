package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	BASE_URL = "https://getpocket.com/v3/"

	PATH_OAUTH_REQUEST = "oauth/request"

	PATH_OAUTH_AUTHZ = "oauth/authorize"
)

const (
	EndpointURL = "https://getpocket.com/v3/get"
)

type PocketResponse struct {
	Status     int                   `json:"status"`
	Complete   int                   `json:"complete"`
	List       map[string]PocketItem `json:"list"`
	Error      interface{}           `json:"error"`
	SearchMeta SearchMeta            `json:"search_meta"`
	Since      int                   `json:"since"`
}

type PocketItem struct {
	ItemID                 string `json:"item_id"`
	ResolvedID             string `json:"resolved_id"`
	GivenURL               string `json:"given_url"`
	GivenTitle             string `json:"given_title"`
	Favorite               string `json:"favorite"`
	Status                 string `json:"status"`
	TimeAdded              string `json:"time_added"`
	TimeUpdated            string `json:"time_updated"`
	TimeRead               string `json:"time_read"`
	TimeFavorited          string `json:"time_favorited"`
	SortID                 int    `json:"sort_id"`
	ResolvedTitle          string `json:"resolved_title"`
	ResolvedURL            string `json:"resolved_url"`
	Excerpt                string `json:"excerpt"`
	IsArticle              string `json:"is_article"`
	IsIndex                string `json:"is_index"`
	HasVideo               string `json:"has_video"`
	HasImage               string `json:"has_image"`
	WordCount              string `json:"word_count"`
	Lang                   string `json:"lang"`
	TopImageURL            string `json:"top_image_url"`
	ListenDurationEstimate int    `json:"listen_duration_estimate"`
}

type SearchMeta struct {
	SearchType string `json:"search_type"`
}

type AuthRequest struct {
	ConsumerKey string `json:"consumer_key"`
	Code        string `json:"code"`
}

func FetchItems(consumerKey string, accessToken string) *PocketResponse {
	params := url.Values{}
	params.Set("consumer_key", consumerKey)
	params.Set("access_token", accessToken)

	resp, err := http.Get(EndpointURL + "?" + params.Encode())
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil
	}

	var pocketResponse PocketResponse
	err = json.Unmarshal(body, &pocketResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return nil
	}

	return &pocketResponse
}

func GetRequestCode(consumerKey string) (string, error) {
	values := url.Values{}
	values.Add("consumer_key", consumerKey)
	values.Add("redirect_uri", "pocketreminder://redirect")

	// リクエストコード取得
	res, err := http.PostForm(
		BASE_URL+PATH_OAUTH_REQUEST,
		values,
	)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	code := strings.Split(string(body), "=")[1]
	fmt.Printf("code: %v", code)

	return code, nil
}

func Authorize(consumerKey string, code string) error {
	// 認証
	authReq := &AuthRequest{
		ConsumerKey: consumerKey,
		Code:        code,
	}

	json, _ := json.Marshal(authReq)
	fmt.Println(string(json))

	res, err := http.Post(
		BASE_URL+PATH_OAUTH_AUTHZ,
		"application/json",
		bytes.NewBuffer(json),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("body: %s", body)

	return nil
}
