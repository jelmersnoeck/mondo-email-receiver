package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-martini/martini"
	"github.com/jelmersnoeck/mondo-email-receiver/gmail"
	"github.com/joho/godotenv"
)

func main() {
	if martini.Env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	g := gmail.NewGmailClient("mondoreceipt@gmail.com")

	m := martini.Classic()

	m.Get("/", func() string {
		return "Hello Mondo crowd!"
	})

	m.Get("/emails/:id", func(params martini.Params) string {
		mail, err := g.Email(params["id"])

		if err != nil {
			return err.Error()
		}

		jsonData, err := json.Marshal(mail)

		if err != nil {
			return err.Error()
		}

		go func(jsonData []byte) {
			url := os.Getenv("WEBHOOK_URL")
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
		}(jsonData)
		return string(jsonData)

		html, err := mail.HTML()

		if err != nil {
			return err.Error()
		}

		return html
	})

	m.Run()
}
