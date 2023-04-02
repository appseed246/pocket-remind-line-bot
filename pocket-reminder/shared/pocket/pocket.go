package pocket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	BASE_URL = "https://getpocket.com/v3/"

	PATH_OAUTH_REQUEST = "oauth/request"
	PATH_OAUTH_AUTHZ   = "oauth/authorize"
	PATH_GET           = "get"
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

func FetchItems(consumerKey string, accessToken string) *PocketResponse {
	params := url.Values{}
	params.Set("consumer_key", consumerKey)
	params.Set("access_token", accessToken)

	resp, err := http.Get(BASE_URL + PATH_GET + "?" + params.Encode())
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
