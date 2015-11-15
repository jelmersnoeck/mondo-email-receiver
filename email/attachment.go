package email

import (
	"encoding/base64"
	"fmt"
)

// Attachment represents an attachment that is sent with an email.
type Attachment struct {
	Body     string `json:"body"`
	MimeType string `json:"mime-type"`
	Filename string `json:"filename"`
}

// HTML will return the value that can be used to put in a HTML attribute.
func (a Attachment) HTML() string {
	data, _ := base64.URLEncoding.DecodeString(a.Body)
	html := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", a.MimeType, html)
}

// Placeholder will return the placeholder that is used in an email to identify
// the attachment.
func (a Attachment) Placeholder() string {
	return fmt.Sprintf("cid:%s", a.Filename)
}
