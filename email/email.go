package email

import "encoding/base64"

// Email represents an email fetched from your gmail account.
type Email struct {
	Subject string `json:"subject"`
	Body    string `json:"body"` // Base64.URLEncoding
	Id      string `json:"id"`
	Sender  string `json:"sender"`
}

// HTML the decoded email body as HTML contents and an error if decoding failed.
func (e Email) HTML() (string, error) {
	data, err := base64.URLEncoding.DecodeString(e.Body)
	return string(data), err
}
