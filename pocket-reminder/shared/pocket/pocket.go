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
	PATH_MODIFY        = "send"
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

type ModifyAction struct {
	Action string `json:"action"`
	ItemId string `json:"item_id"`
	Time   string `json:"time"`
}

type ModifyResponse struct {
	// ActionResults []bool `json:"action_results"`
	ActionErrors []bool `json:"action_errors"`
	Status       uint8  `json:"status"`
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

func ModifyItem(consumerKey string, accessToken string, actions *[]ModifyAction) (*ModifyResponse, error) {
	// actionsのJsonArray文字列への変換
	jsonData, err := json.Marshal(*actions)
	if err != nil {
		fmt.Println("Error marshalling JSON: ", err)
		return nil, err
	}

	params := url.Values{}
	params.Set("consumer_key", consumerKey)
	params.Set("access_token", accessToken)
	params.Set("actions", string(jsonData))

	resp, err := http.Get(BASE_URL + PATH_MODIFY + "?" + params.Encode())
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}
	fmt.Println(string(body))

	var r ModifyResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return nil, err
	}

	return &r, nil
}
