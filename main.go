package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

//type Event struct {
//	ID        string    `json:"id"`
//	Type      string    `json:"type"`
//	CreatedAt time.Time `json:"created_at"`
//	Repo      struct {
//		Name string `json:"name"`
//	} `json:"repo"`
//}

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		PushID  int64 `json:"push_id"`
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
		Action string `json:"action"`
		Issue  struct {
			Title string `json:"title"`
		} `json:"issue"`
		Forkee struct {
			FullName string `json:"full_name"`
		} `json:"forkee"`
	} `json:"payload"`
}

func main() {
	username := os.Args[1]
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error ferching")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		return
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			commitCount := len(event.Payload.Commits)
			fmt.Printf("Pushed %d commit(s) to %s\n", commitCount, event.Repo.Name)
		case "IssuesEvent":
			fmt.Printf("Opened a new issue in %s: %s\n", event.Repo.Name, event.Payload.Issue.Title)
		case "WatchEvent":
			fmt.Printf("Starred %s\n", event.Repo.Name)
		case "ForkEvent":
			fmt.Printf("Forked %s to %s\n", event.Repo.Name, event.Payload.Forkee.FullName)
		case "CreateEvent":
			fmt.Printf("Created a new %s in %s\n", event.Payload.Action, event.Repo.Name)
		case "PullRequestEvent":
			fmt.Printf("%s a pull request in %s\n", strings.Title(event.Payload.Action), event.Repo.Name)
		case "IssueCommentEvent":
			fmt.Printf("Commented on an issue in %s\n", event.Repo.Name)
		default:
			fmt.Printf("%s event occurred in %s\n", event.Type, event.Repo.Name)
		}
	}
}
