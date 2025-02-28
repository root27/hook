package main

import (
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//github username
	//github token
	//webhook url
	//webhook secret

	fmt.Println("Welcome to the Github Webhook CLI")

	var username string

	fmt.Print(color.YellowString("Enter your Github username (required): "))

	fmt.Scanln(&username)

	var token string

	fmt.Print(color.YellowString("Enter your Github token (required): "))
	fmt.Scanln(&token)

	var url string

	fmt.Print(color.YellowString("Enter the webhook url (required): "))
	fmt.Scanln(&url)

	var secret string

	fmt.Print(color.YellowString("Enter the webhook secret (suggested): "))
	fmt.Scanln(&secret)

	if url == "" || token == "" || username == "" {
		fmt.Println(color.RedString("All fields are required"))
		return
	}

	cli := Cli{
		Username:      username,
		Token:         token,
		WebhookUrl:    url,
		WebhookSecret: secret,
		HttpClient:    &http.Client{},
	}

	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopChan
		fmt.Println(color.RedString("Cli Stopped"))
		os.Exit(0)
	}()

	cli.HandleRepos()

}
