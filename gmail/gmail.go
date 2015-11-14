package gmail

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jelmersnoeck/mondo-email-receiver/email"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tok, _ := tokenFromFile("credentials.json")
	return config.Client(ctx, tok)
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// GmailClient represent a connection with Gmail
type GmailClient struct {
	srv   *gmail.Service
	email string
}

func NewGmailClient(email string) GmailClient {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

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
	return ""
}
