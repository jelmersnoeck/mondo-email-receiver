package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-martini/martini"
	"github.com/jelmersnoeck/mondo-email-receiver/gmail"
)

func main() {
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
			fmt.Println(string(jsonData))
		}(jsonData)

		html, err := mail.HTML()

		if err != nil {
			return err.Error()
		}

		return html
	})

	m.Run()
}
