package convert

import (
	"fmt"
	"net/url"
	"pocket-reminder/shared/pocket"
	"strings"
	"unicode/utf8"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func CreateCarouselMessage(response *pocket.PocketResponse) []*linebot.CarouselColumn {
	var columns []*linebot.CarouselColumn

	counter := 0
	for _, item := range response.List {
		// MessagingAPIの制約でカルーセルアイテムの上限を越さないようにする
		if counter >= 10 {
			break
		}
		column := createCarouselColumnFromPocketItem(item)
		columns = append(columns, column)
		counter++
	}

	return columns
}

func createCarouselColumnFromPocketItem(item pocket.PocketItem) *linebot.CarouselColumn {
	// MessagingAPIの制約
	truncatedTitle := truncateString(item.GivenTitle, 40)
	truncatedText := truncateString(item.ResolvedTitle, 60)

	fmt.Println(item.TopImageURL)
	imageURL := item.TopImageURL
	if imageURL == "" || strings.HasPrefix(imageURL, "http://") {
		imageURL = "https://sozaino.site/wp-content/uploads/2023/01/sen-simple53.png"
	}

	// タイトルが存在しない場合、デフォルト値を設定する
	if truncatedTitle == "" {
		truncatedTitle = "No title"
	}

	// 外部ブラウザで開くよう末尾にリクエストパラメータを追加
	parsedURL, _ := url.Parse(item.GivenURL)
	query := parsedURL.Query()
	query.Set("openExternalBrowser", "1")
	parsedURL.RawQuery = query.Encode()
	modifiedURL := parsedURL.String()

	actions := []linebot.TemplateAction{
		linebot.NewURIAction("ブラウザで開く", modifiedURL),
		linebot.NewMessageAction("アーカイブ", fmt.Sprintf("/archive %s", item.ItemID)),
	}
	return linebot.NewCarouselColumn(imageURL, truncatedTitle, truncatedText, actions...)
}

func truncateString(s string, maxLength int) string {
	suffix := "..."
	if utf8.RuneCountInString(s) <= maxLength-len(suffix) {
		return s
	}
	return string([]rune(s)[:maxLength-len(suffix)]) + suffix
}

func PrintCarouselColumns(columns []*linebot.CarouselColumn) {
	for index, column := range columns {
		actionsText := ""

		for actionIndex, action := range column.Actions {
			switch act := action.(type) {
			case *linebot.URIAction:
				actionsText += fmt.Sprintf("Action %d (URIAction): Label=%s, URI=%s; ", actionIndex+1, act.Label, act.URI)
			default:
				actionsText += fmt.Sprintf("Action %d (Unknown); ", actionIndex+1)
			}
		}

		fmt.Printf("Column %d: ThumbnailURL=%s, Title=%s, Text=%s, %s\n", index+1, column.ThumbnailImageURL, column.Title, column.Text, actionsText)
	}
}
