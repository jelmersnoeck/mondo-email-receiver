package email

import "encoding/base64"

// Email represents an email fetched from your gmail account.
type Email struct {
	Subject string
	// Base64 URLEncoded content
	Body   string
	Id     string
	Sender string
}

// HTML the decoded email body as HTML contents and an error if decoding failed.
func (e Email) HTML() (string, error) {
	data, err := base64.URLEncoding.DecodeString(e.Body)
	return string(data), err
}
