package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/root27/hook/internal"
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

	var number_events int

	fmt.Print(color.YellowString("Enter the number of events that you want (default:1): "))
	fmt.Scanln(&number_events)

	if number_events == 0 {
		number_events = 1
	}

	var events []string

	for i := 0; i < number_events; i++ {

		var event string
		fmt.Print(color.RedString("Enter the event name (default:push): "))
		fmt.Scanln(&event)

		if event == "" {

			event = "push"

		}

		events = append(events, event)

	}

	if url == "" || token == "" || username == "" {
		fmt.Println(color.RedString("All fields are required"))
		return
	}

	cli := internal.Cli{
		Username:      username,
		Token:         token,
		WebhookUrl:    url,
		WebhookSecret: secret,
		Events:        events,
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
