package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile("credentials.json")
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

func main() {
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

	call := srv.Users.Messages.Get("mondoreceipt@gmail.com", "1510670a190c8455")

	res, _ := call.Format("full").Do()

	data := ""
	showParts(res.Payload.Parts)

	fmt.Println(data)
}

func showParts(parts []*gmail.MessagePart) {
	for _, part := range parts {
		if len(part.Parts) > 0 {
			showParts(part.Parts)
		} else {
			if part.MimeType == "text/html" {
				data, _ := base64.URLEncoding.DecodeString(part.Body.Data)
				fmt.Println(string(data))
			}
		}
	}
}
