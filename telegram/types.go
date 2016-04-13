package telegram

type User struct {
	Id         int64  `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Username   string `json:"username"`
}

type GroupChat struct {
}

type Message struct {
	Message_id            int64
	From                  User
	Date                  int64
	Chat                  User
	Forward_from          User
	Forward_date          int64
	Reply_to_message      *Message
	Text                  string
	Audio                 Audio
	Document              Document
	Photo                 []PhotoSize
	Sticker               Sticker
	Video                 Video
	Voice                 Voice
	Caption               string
	Contact               Contact
	Location              Location
	New_chat_participant  User
	Left_chat_participant User
	New_chat_title        string
	New_chat_photo        []PhotoSize
	Delete_chat_photo     bool
	Group_chat_created    bool
}

type PhotoSize struct {
}

type Audio struct {
}

type Document struct {
}

type Sticker struct {
}

type Video struct {
}

type Voice struct {
}

type Contact struct {
}

type Location struct {
}

type UserProfilePhoto struct {
}

type ReplyKeyboardMarkup struct {
}

type ReplyKeyboardHide struct {
}

type ForceReply struct {
}

type InlineQuery struct {
	Id     string `json:"id"`
	From   User   `json:"from"`
	Query  string `json:"query"`
	Offset string `json:"offset"`
}

type ChosenInlineQuery struct {
	Result_id         string `json:"result_id"`
	From              User   `json:"from"`
	Inline_message_id string `json:"inline_message_id"`
	Query             string `json:"query"`
}

type AnswerInlineQuery struct {
	Inline_query_id string
	Results         []InlineQueryResultPhoto
}

type InlineQueryResultPhoto struct {
	Type        string `json:"type"`
	Id          string `json:"id"`
	Photo_url   string `json:"photo_url"`
	Thumb_url   string `json:"thumb_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Caption     string `json:"caption"`
}

type Update struct {
	Update_id            int64             `json:"update_id"`
	Message              Message           `json:"message"`
	Inline_query         InlineQuery       `json:"inline_query"`
	Chosen_inline_result ChosenInlineQuery `json:"chosen_inline_result"`
}
