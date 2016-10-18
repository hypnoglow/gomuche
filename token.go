package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"golang.org/x/oauth2"
)

const tokenFilename = "$HOME/.gomuche/token.json"

// NewTokenFromFile reads token from file and returns it.
func NewTokenFromFile() *oauth2.Token {
	filename := os.ExpandEnv(tokenFilename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	token := new(oauth2.Token)

	err = json.Unmarshal(bytes, &token)
	if err != nil {
		log.Fatalln(err)
	}

	return token
}

// NewTokenFromCode retrieves a new token from Google using code and returns it.
func NewTokenFromCode(conf *oauth2.Config, code string) *oauth2.Token {
	ctx := context.Background()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatalln(err)
	}

	SaveToken(token)

	return token
}

// SaveToken saves token to file.
func SaveToken(token *oauth2.Token) {
	bytes, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	filename := os.ExpandEnv(tokenFilename)
	err = os.MkdirAll(path.Dir(filename), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(filename, bytes, 0755)
	if err != nil {
		log.Fatalln(err)
	}
}
