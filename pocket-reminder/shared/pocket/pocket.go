package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	BASE_URL = "https://getpocket.com/"

	PATH_OAUTH_REQUEST = "v3/oauth/request"
	PATH_OAUTH_AUTHZ   = "auth/authorize"
	PATH_OAUTH_TOKEN   = "v3/oauth/authorize"
	PATH_GET           = "v3/get"
	PATH_MODIFY        = "v3/send"
)

type Client struct {
	BaseURL     string
	ConsumerKey string
	AccessToken string
}

type GetResponse struct {
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

func New(consumerKey string, accessToken string) (*Client, error) {
	if consumerKey == "" {
		return nil, errors.New("consumer key must not be empty")
	}

	if accessToken == "" {
		return nil, errors.New("access token must not be empty")
	}

	return &Client{
		BaseURL:     BASE_URL,
		ConsumerKey: consumerKey,
		AccessToken: accessToken,
	}, nil
}

func GetAuthorizationCode(consumerKey string, redirectURL string) (string, error) {
	params := &struct {
		CosumerKey  string `json:"consumer_key"`
		RedirectURL string `json:"redirect_uri"`
	}{
		CosumerKey:  consumerKey,
		RedirectURL: redirectURL,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return "", nil
	}

	resp, err := http.Post(BASE_URL+PATH_OAUTH_REQUEST, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return "", err
	}

	f, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Printf("Error parse formData: %v\n", err)
		return "", nil
	}

	return f.Get("code"), nil
}

func GetAccessToken(consumerKey string, authorizationCode string) (token string, userName string, e error) {
	params := &struct {
		CosumerKey string `json:"consumer_key"`
		Code       string `json:"code"`
	}{
		CosumerKey: consumerKey,
		Code:       authorizationCode,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return "", "", nil
	}

	resp, err := http.Post(BASE_URL+PATH_OAUTH_TOKEN, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return "", "", err
	}

	f, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Printf("Error parse formData: %v\n", err)
		return "", "", nil
	}

	return f.Get("access_token"), f.Get("username"), nil
}

func (c *Client) FetchItems() *GetResponse {
	params := url.Values{}
	params.Set("consumer_key", c.ConsumerKey)
	params.Set("access_token", c.AccessToken)

	resp, err := http.Get(c.BaseURL + PATH_GET + "?" + params.Encode())
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
	fmt.Println(string(body))

	var pocketResponse GetResponse
	err = json.Unmarshal(body, &pocketResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return nil
	}

	return &pocketResponse
}

func (c *Client) ModifyItem(actions *[]ModifyAction) (*ModifyResponse, error) {
	// actionsのJsonArray文字列への変換
	jsonData, err := json.Marshal(*actions)
	if err != nil {
		fmt.Println("Error marshalling JSON: ", err)
		return nil, err
	}

	params := url.Values{}
	params.Set("consumer_key", c.ConsumerKey)
	params.Set("access_token", c.AccessToken)
	params.Set("actions", string(jsonData))

	resp, err := http.Get(c.BaseURL + PATH_MODIFY + "?" + params.Encode())
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
