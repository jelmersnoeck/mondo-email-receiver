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

// GmailClient represent a connection with Gmail
type GmailClient struct {
	srv   *gmail.Service
	email string
}

func NewGmailClient(email string) GmailClient {
	ctx := context.Background()

	client := getClient(ctx, getConfig())

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	return GmailClient{srv, email}
}

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

	return email, nil
}

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

func getMessageSender(headers []*gmail.MessagePartHeader) string {
	return getMessageHeader(headers, "From")
}

func getMessageSubject(headers []*gmail.MessagePartHeader) string {
	return getMessageHeader(headers, "Subject")
}

func getMessageHeader(headers []*gmail.MessagePartHeader, wanted string) string {
	for _, header := range headers {
		if header.Name == wanted {
			return header.Value
		}

	}

	return ""
}
