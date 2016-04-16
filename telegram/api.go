package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const baseurl string = "https://api.telegram.org/bot"

type GetUpdatesResponse struct {
	Ok     bool
	Result []Update
}

func StartFetchUpdates(updateChannel *chan []Update) {

	var since int64 = 0
	defer close(*updateChannel)

	for {
		updates := GetUpdates(since)
		if len(updates) > 0 {
			since = updates[len(updates)-1].Update_id + 1
		}
		*updateChannel <- updates
		time.Sleep(1 * time.Second)
	}

}

func SendMessage(chatId int64, text string) {
	url := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=%s", baseurl, Token, chatId, text)
	_, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
}

func AnswerQuery(queryId string, results []InlineQueryResultPhoto) {
	json, err := json.Marshal(results)

	if err != nil {
		fmt.Println(err)
		return
	}
	url := fmt.Sprintf("%s%s/answerInlineQuery?inline_query_id=%s&results=%s", baseurl, Token, queryId, json)
	_, err = http.Get(url)

	if err != nil {
		fmt.Println(err)
	}
}

func GetUpdates(offset int64) []Update {
	url := baseurl + Token + "/getUpdates"
	if offset != 0 {
		url += fmt.Sprintf("?offset=%d", offset)
	}

	response, err := http.Get(url)

	if err != nil {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil
	}

	var result GetUpdatesResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return result.Result
}
