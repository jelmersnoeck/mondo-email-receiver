package main

import (
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
		email, err := g.Email(params["id"])

		if err != nil {
			return err.Error()
		}

		html, err := email.HTML()

		if err != nil {
			return err.Error()
		}

		return html
	})

	m.Run()
}
