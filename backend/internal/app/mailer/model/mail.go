package model

type Mail struct {
	To      string
	Subject string
	Body    string
}

type MailConfig struct {
	APIKey  string
	BaseURL string
}
