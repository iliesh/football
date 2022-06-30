package telegram

type webHookReqBodyT struct {
	UpdateID      int            `json:"update_id"`
	Message       messageT       `json:"message,omitempty"`
	CallBackQuery callBackQueryT `json:"callback_query,omitempty"`
}

// messageT Struct
type messageT struct {
	MessageID            int64  `json:"message_id"`
	From                 userT  `json:"from"`
	ForwardFrom          userT  `json:"forward_from"`
	Chat                 chatT  `json:"chat"`
	SenderChat           chatT  `json:"sender_chat"`
	ForwardFromChat      chatT  `json:"forward_from_chat"`
	ForwardFromMessageID int64  `json:"forward_from_message_id"`
	ForwardSignature     string `json:"forward_signature"`
	ForwardSenderName    string `json:"forward_sender_name"`
	ForwardDate          int64  `json:"forward_date"`
	Date                 int64  `json:"date"`
	Text                 string `json:"text"`
	Entities             []struct {
		Offset int    `json:"offset"`
		Length int    `json:"length"`
		Type   string `json:"type"`
	} `json:"entities"`
	ViaBot          userT  `json:"via_bot"`
	EditDate        int64  `json:"edit_date"`
	MediaGroupID    string `json:"media_group_id"`
	AuthorSignature string `json:"author_signature"`
	Voice           voiceT `json:"voice"`
	Caption         string `json:"caption"`
}

// callBackQueryT Struct
type callBackQueryT struct {
	ID              string   `json:"id"`
	From            userT    `json:"from"`
	Message         messageT `json:"message,omitempty"`
	InlineMessageID string   `json:"inline_message_id,omitempty"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
}

// userT Struct
type userT struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// chatT struct
type chatT struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// voiceT struct
type voiceT struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

// sendMessageReqBodyT Create a struct to conform to the JSON body
// of the send message request
// https://core.telegram.org/bots/api#sendmessage
type sendMessageReqBodyT struct {
	ChatID      int64                 `json:"chat_id,omitempty"`
	Text        string                `json:"text,omitempty"`
	ParseMode   string                `json:"parse_mode,omitempty"`
	ReplyMarkup inlineKeyboardMarkupT `json:"reply_markup,omitempty"`
}

type inlineKeyboardMarkupT struct {
	InlineKeyboard [][]inlineKeyboardButtonT `json:"inline_keyboard,omitempty"`
}

type inlineKeyboardButtonT struct {
	Text         string `json:"text,omitempty"`
	URL          string `json:"url,omitempty"`
	CallBackData string `json:"callback_data,omitempty"`
}
