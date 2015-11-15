// Package gmail implements a way to fetch emails from Gmail.

package gmail

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jelmersnoeck/mondo-email-receiver/email"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

// GmailClient represent a connection with Gmail
type GmailClient struct {
	srv   *gmail.Service
	email string
}

// NewGmailClient will set up a new GmailClient for the specified email address.
func NewGmailClient(email string) GmailClient {
	ctx := context.Background()

	client := getClient(ctx, getConfig())

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	return GmailClient{srv, email}
}

// Email will return an email type that contains the HTML body, subject, sender
// and related attachments of an email specified by it's ID.
func (c *GmailClient) Email(id string) (email.Email, error) {
	call := c.srv.Users.Messages.Get(c.email, id)
	res, err := call.Format("full").Do()

	email := email.Email{}

	if err != nil {
		return email, err
	}

	email.Id = id
	email.Body = getMessageBody(res.Payload.Parts)
	email.Sender = getMessageSender(res.Payload.Headers)
	email.Subject = getMessageSubject(res.Payload.Headers)
	email.Attachments = c.getMessageAttachments(
		res.Payload.Parts,
		id,
		c.srv.Users.Messages.Attachments,
	)

	return email, nil
}

// getMessageBody finds the HTML body of an email.
func getMessageBody(parts []*gmail.MessagePart) string {
	for _, part := range parts {
		if len(part.Parts) > 0 {
			return getMessageBody(part.Parts)
		} else {
			if part.MimeType == "text/html" {
				return part.Body.Data
			}
		}
	}

	return ""
}

// getMessageAttachments goes through the message parts for a specific message
// to find all the image attachments.
func (c *GmailClient) getMessageAttachments(
	parts []*gmail.MessagePart,
	messageId string,
	s *gmail.UsersMessagesAttachmentsService) []email.Attachment {

	attachments := make([]email.Attachment, 0)

	for _, part := range parts {
		if len(part.Parts) == 0 {
			if isImageAttachment(part.MimeType) {
				gmailAttachment := s.Get(
					c.email,
					messageId,
					part.Body.AttachmentId,
				)
				body, err := gmailAttachment.Do()

				if err != nil {
					continue
				}

				attachment := email.Attachment{
					Body:     body.Data,
					MimeType: part.MimeType,
					Filename: part.Filename,
				}
				attachments = append(attachments, attachment)
			}
		}
	}

	return attachments
}

var imageMimeTypes = []string{
	"image/png",
	"image/jpg",
	"image/jpeg",
	"image/gif",
}

// isImageAttachment validates a mime type to see if it is an image MimeType.
func isImageAttachment(mime string) bool {
	for _, tp := range imageMimeTypes {
		if tp == mime {
			return true
		}
	}

	return false
}

// getMessageSender goes through the headers to find the From header.
func getMessageSender(headers []*gmail.MessagePartHeader) string {
	return getMessageHeader(headers, "From")
}

// getMessageSubject goes through the headers to find the Subject header.
func getMessageSubject(headers []*gmail.MessagePartHeader) string {
	return getMessageHeader(headers, "Subject")
}

// getMessageHeader goes through a list of headers and returns the header where
// the name matches the one we want.
func getMessageHeader(headers []*gmail.MessagePartHeader, wanted string) string {
	for _, header := range headers {
		if header.Name == wanted {
			return header.Value
		}

	}

	return ""
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	return config.Client(ctx, oauth2Token())
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func oauth2Token() *oauth2.Token {
	tm, _ := time.Parse("2006-Jan-02", "2015-Nov-13")
	return &oauth2.Token{
		AccessToken:  os.Getenv("GOOGLE_ACCESS_TOKEN"),
		TokenType:    "Bearer",
		RefreshToken: os.Getenv("GOOGLE_REFRESH_TOKEN"),
		Expiry:       tm,
	}
}

func getConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}
