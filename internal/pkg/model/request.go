package model

type CreateChatRequest struct {
	Title string `json:"title"`
}

type SendMessageRequest struct {
	Text string `json:"text"`
}

type GetChatRequest struct {
	Limit int `json:"limit"`
}
