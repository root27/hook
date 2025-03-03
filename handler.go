package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Cli struct {
	Username      string
	Token         string
	WebhookUrl    string
	WebhookSecret string
	Events        []string
	HttpClient    *http.Client
}

type Webhook any

type Repository any

func (c *Cli) HandleRepos() {

	//1 -> Get all repos
	//2 -> Check each repo if webhook is present
	//3-> If webhook is not present, create webhook

	repos, err := c.GetRepos()

	if err != nil {

		fmt.Println("Error fetching repos: ", err)
	}

	fmt.Println("Checking webhooks for repositories:")
	for _, repo := range repos {
		webhookExists, err := c.CheckWebHook(repo.(map[string]any)["name"])
		if err != nil {
			fmt.Printf("Error checking webhook for %s: %v\n", repo.(map[string]any)["name"], err)
			continue
		}
		if webhookExists {
			fmt.Printf("Webhook exists for %s\n", repo.(map[string]any)["name"])
		} else {
			fmt.Printf("No webhook found for %s\n", repo.(map[string]any)["name"])

			err := c.CreateWebHook(repo)

			if err != nil {

				fmt.Printf("Error creating webhook for %s: %v\n", repo.(map[string]any)["name"], err)

			}

		}
	}

	fmt.Println("Webhook created for all repositories")

}

func (c *Cli) CreateWebHook(repo Repository) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", c.Username, repo.(map[string]any)["name"])
	payload := map[string]any{
		"name":   "web",
		"active": true,
		"events": []string{"push"},
		"config": map[string]any{
			"url":          c.WebhookUrl,
			"content_type": "json",
			"secret":       c.WebhookSecret,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create webhook: %s", resp.Status)
	}
	return nil
}

func (c *Cli) CheckWebHook(repo Repository) (bool, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", c.Username, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "token "+c.Token)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to fetch webhooks: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var webhooks []Webhook
	if err := json.Unmarshal(body, &webhooks); err != nil {
		return false, err
	}

	return len(webhooks) > 0, nil
}

func (c *Cli) GetRepos() ([]Repository, error) {

	var allRepos []Repository
	page := 1

	for {
		url := fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&per_page=10", c.Username, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch repos: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var repos []Repository
		if err := json.Unmarshal(body, &repos); err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		// Check for pagination
		linkHeader := resp.Header.Get("Link")
		if !strings.Contains(linkHeader, "rel=\"next\"") {
			break // No more pages
		}
		page++
	}

	return allRepos, nil

}
