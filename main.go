package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	EnvSlackWebhook  = "SLACK_WEBHOOK"
	EnvSlackIcon     = "SLACK_ICON"
	EnvSlackChannel  = "SLACK_CHANNEL"
	EnvPrimarySlackChannel  = "PRIMARY_SLACK_CHANNEL"
	EnvSlackTitle    = "SLACK_TITLE"
	EnvSlackMessage  = "SLACK_MESSAGE"
	EnvSlackColor    = "SLACK_COLOR"
	EnvSlackUserName = "SLACK_USERNAME"
	EnvGithubActor   = "GITHUB_ACTOR"
	EnvSiteName      = "SITE_NAME"
	EnvHostName      = "HOST_NAME"
	EnvDepolyPath    = "DEPLOY_PATH"
)

type Webhook struct {
	Text        string       `json:"text,omitempty"`
	UserName    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links"`
	Attachments []Attachment `json:"attachments,omitmepty"`
}

type Attachment struct {
	Fallback   string  `json:"fallback"`
	Pretext    string  `json:"pretext,omitempty"`
	Color      string  `json:"color,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"`
	AuthorIcon string  `json:"author_icon,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
	
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

func main() {
	endpoint := os.Getenv(EnvSlackWebhook)
	if endpoint == "" {
		fmt.Fprintln(os.Stderr, "URL is required")
		os.Exit(1)
	}
	text := os.Getenv(EnvSlackMessage)
	if text == "" {
		fmt.Fprintln(os.Stderr, "Message is required")
		os.Exit(1)
	}
	var branch string
	var ref string
	if len(strings.SplitN(os.Getenv("GITHUB_REF"), "/",100)) > 2 {
		branch = strings.SplitN((os.Getenv("GITHUB_REF")), "/", 100)[2]
		if strings.SplitN(os.Getenv("GITHUB_REF"), "/",100)[1] == "heads" {
			ref="Branch"
		} else if strings.SplitN(os.Getenv("GITHUB_REF"), "/",100)[1] == "tags" {
			ref="Tag"
		} else {
			ref="Ref"
		}
	}	else {
		branch = os.Getenv("GITHUB_REF")
		ref="Ref"
	} 
	fields:= []Field{
		{
			Title: ref,
			Value: branch,
			Short: true,
		},                {
			Title: "Event",
			Value: os.Getenv("GITHUB_EVENT_NAME"),
			Short: true,
		},
		{
			Title: "Actions URL",
			Value: "https://github.com/" + os.Getenv("GITHUB_REPOSITORY") + "/commit/" + os.Getenv("GITHUB_SHA") + "/checks",
			Short: false,
		},
		{
			Title: os.Getenv(EnvSlackTitle),
			Value: envOr(EnvSlackMessage, "EOM"),
			Short: false,
		},
	}

	hostName := os.Getenv(EnvHostName)
	if hostName != "" {
		newfields:= []Field{
			{
				Title: os.Getenv("SITE_TITLE"),
				Value: os.Getenv(EnvSiteName),
				Short: true,
			},
			{
				Title: os.Getenv("HOST_TITLE"),
				Value: os.Getenv(EnvHostName),
				Short: true,
			},
		}
		fields = append(newfields, fields...)
	}

	msg := Webhook{
		UserName: os.Getenv(EnvSlackUserName),
		IconURL:  os.Getenv(EnvSlackIcon),
		Channel:  os.Getenv(EnvSlackChannel),
		Attachments: []Attachment{
			{
				Fallback: envOr(EnvSlackMessage, "GITHUB_ACTION=" + os.Getenv("GITHUB_ACTION") + " \n GITHUB_ACTOR=" + os.Getenv("GITHUB_ACTOR") + " \n GITHUB_EVENT_NAME=" + os.Getenv("GITHUB_EVENT_NAME") + " \n GITHUB_REF=" + os.Getenv("GITHUB_REF") + " \n GITHUB_REPOSITORY=" + os.Getenv("GITHUB_REPOSITORY") + " \n GITHUB_WORKFLOW=" + os.Getenv("GITHUB_WORKFLOW")),
				Color:      envOr(EnvSlackColor, "good"),
				AuthorName: envOr(EnvGithubActor, ""),
				AuthorLink: "http://github.com/" + os.Getenv(EnvGithubActor),
				AuthorIcon: "http://github.com/" + os.Getenv(EnvGithubActor) + ".png?size=32",
				Fields: fields,
			},
		},
	}

//double send message to primary channel
	if err := send(endpoint, msg); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err)
		os.Exit(2)
	}
	if EnvPrimarySlackChannel != "" {
		priMsg := Webhook{
			UserName: os.Getenv(EnvSlackUserName),
			IconURL:  os.Getenv(EnvSlackIcon),
			Channel:  os.Getenv(EnvPrimarySlackChannel),
			Attachments: []Attachment{
				{
					Fallback: envOr(EnvSlackMessage, "GITHUB_ACTION=" + os.Getenv("GITHUB_ACTION") + " \n GITHUB_ACTOR=" + os.Getenv("GITHUB_ACTOR") + " \n GITHUB_EVENT_NAME=" + os.Getenv("GITHUB_EVENT_NAME") + " \n GITHUB_REF=" + os.Getenv("GITHUB_REF") + " \n GITHUB_REPOSITORY=" + os.Getenv("GITHUB_REPOSITORY") + " \n GITHUB_WORKFLOW=" + os.Getenv("GITHUB_WORKFLOW")),
					Color:      envOr(EnvSlackColor, "good"),
					AuthorName: envOr(EnvGithubActor, ""),
					AuthorLink: "http://github.com/" + os.Getenv(EnvGithubActor),
					AuthorIcon: "http://github.com/" + os.Getenv(EnvGithubActor) + ".png?size=32",
					Fields: fields,
				},
			},
		}


	 	if err := send(endpoint, priMsg); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err)
			os.Exit(2)
		}
	}
}


func envOr(name, def string) string {
	if d, ok := os.LookupEnv(name); ok {
		return d
	}
	return def
}

func send(endpoint string, msg Webhook) error {
	enc, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(enc)
	res, err := http.Post(endpoint, "application/json", b)
	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		return fmt.Errorf("Error on message: %s\n", res.Status)
	}
	fmt.Println(res.Status)
	return nil
}