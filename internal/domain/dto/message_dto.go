package dto

type EmailMessage struct {
	To       string `json:"to"`
	Template string `json:"template"`
	Subject  string `json:"subject"`
	Data     any    `json:"data"`
}
