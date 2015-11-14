package main

import (
	"fmt"

	"github.com/jelmersnoeck/mondo-email-receiver/gmail"
)

func main() {
	g := gmail.NewGmailClient("mondoreceipt@gmail.com")

	email, _ := g.Email("1510670a190c8455")
	html, _ := email.HTML()
	fmt.Println(html)
}
