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

func (a Attachment) HTML() string {
	data, _ := base64.URLEncoding.DecodeString(a.Body)
	html := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", a.MimeType, html)
}

func (a Attachment) Placeholder() string {
	return fmt.Sprintf("cid:%s", a.Filename)
}
