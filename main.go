package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
	"golang.org/x/oauth2"
	"os"
	"io"
)

const (
	usage = `gomuche - Google Mail Unread count checker.

Usage:
  gomuche auth
  gomuche check [-v] [-c <code>]
  gomuche -h | --help
  gomuche -V | --version

Check options:
  -c --code=<authCode>    Auth code, which can be obtained
                          through 'gomuche auth' command.
  -v --verbose            Verbose output. This shows errors
                          instead of just silently exiting with error code.

Other options:
  -h --help               Show this helpful info.
  -V --version            Show version.
`
	version = "0.1.0"
)

const (
	mailFeedURL = "https://mail.google.com/mail/feed/atom"
	authURL     = "https://accounts.google.com/o/oauth2/auth"
	tokenURL    = "https://www.googleapis.com/oauth2/v4/token"
	redirectURL = "urn:ietf:wg:oauth:2.0:oob"
)

// Feed represents Google Mail Atom Feed.
type Feed struct {
	Title     string     `xml:"title"`
	Tagline   string     `xml:"tagline"`
	Fullcount int        `xml:"fullcount"`
	Modified  *time.Time `xml:"modified"`
}

func main() {
	logFile := getLogFile(true)
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	log.Println("Start.")

	args, err := docopt.Parse(usage, nil, true, "gomuche "+version, false)
	if err != nil {
		log.Fatalln("Error parsing arguments:", err)
	}
	log.Printf("Arguments: %v\n", args)

	isVerbose := parseVerbose(args)
	if isVerbose || args["auth"] == true {
		log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	}

	cfg := NewConfigFromFile()
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		log.Fatalln("Client ID and secret are not specified.")
	}

	oauth2conf := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{mailFeedURL},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		RedirectURL: redirectURL,
	}

	switch {
	case args["auth"] == true:
		authAction(oauth2conf)

	case args["check"] == true:
		code := strings.TrimSpace(parseCode(args))
		checkAction(oauth2conf, code)

	default:
		log.Fatalln("Action is not defined.")
	}
}

func parseCode(args map[string]interface{}) string {
	if args["--code"] == nil {
		return ""
	}

	return args["--code"].(string)
}

func parseVerbose(args map[string]interface{}) bool {
	if args["--verbose"] == nil {
		return false
	}

	return args["--verbose"].(bool)
}

func checkAction(conf *oauth2.Config, code string) {
	var token *oauth2.Token

	if code != "" {
		token = NewTokenFromCode(conf, code)
	} else {
		token = NewTokenFromFile()
	}

	tokenSource := conf.TokenSource(oauth2.NoContext, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalln(err)
	}

	if newToken.AccessToken != token.AccessToken {
		SaveToken(newToken)
		log.Println("Saved new token:", newToken.AccessToken)
	}

	client := oauth2.NewClient(oauth2.NoContext, tokenSource)
	resp, err := client.Get(mailFeedURL)
	if err != nil {
		log.Fatalln("Error fetching mail feed:", err)
	}

	defer resp.Body.Close()

	feed := new(Feed)
	err = xml.NewDecoder(resp.Body).Decode(feed)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(feed.Fullcount)
}

func authAction(conf *oauth2.Config) {
	params := conf.AuthCodeURL("state")
	fmt.Printf("Visit the URL for the auth dialog:\n%v\n", params)
}