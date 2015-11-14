package main

import (
	"fmt"

	"github.com/jelmersnoeck/gmail/gmail"
)

func main() {
	g := gmail.NewGmailClient("mondoreceipt@gmail.com")

	email, _ := g.Email("1510670a190c8455")
	fmt.Println(email.HTML())
}
