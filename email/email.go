package email

import (
	"encoding/base64"
	"strings"
)

// Email represents an email fetched from your gmail account.
type Email struct {
	Subject     string       `json:"subject"`
	Body        string       `json:"body"` // Base64.URLEncoding
	Id          string       `json:"id"`
	Sender      string       `json:"sender"`
	Attachments []Attachment `json:"attachments"`
}

// HTML the decoded email body as HTML contents and an error if decoding failed.
func (e Email) HTML() (html string, err error) {
	data, err := base64.URLEncoding.DecodeString(e.Body)

	if err == nil && len(e.Attachments) > 0 {
		html = string(data)

		if len(e.Attachments) > 0 {
			html = replaceAttachments(html, e.Attachments)
		}
	}

	return html, err
}

func replaceAttachments(html string, attachments []Attachment) string {
	for _, attachment := range attachments {
		html = strings.Replace(
			html,
			attachment.Placeholder(),
			attachment.HTML(),
			-1,
		)
	}

	return html
}
