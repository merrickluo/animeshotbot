package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const TOKEN string = ""
const BASE_URL string = "https://api.telegram.org/bot"

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
	url := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=%s", BASE_URL, TOKEN, chatId, text)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func AnswerQuery(queryId string, results []InlineQueryResultPhoto) {
	json, err := json.Marshal(results)

	if err != nil {
		fmt.Println(err)
		return
	}
	url := fmt.Sprintf("%s%s/answerInlineQuery?inline_query_id=%s&results=%s", BASE_URL, TOKEN, queryId, json)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func GetUpdates(offset int64) []Update {
	url := BASE_URL + TOKEN + "/getUpdates"
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
	var str = fmt.Sprintf("%s", body)
	println(str)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return result.Result
}
