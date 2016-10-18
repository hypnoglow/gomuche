# gomuche

gomuche stands for Google Mail unread count checker.

## Motivation

This tool in fact does this:

    curl -su $GMAIL_USER:$GMAIL_PASS https://mail.google.com/mail/feed/atom
    
But instead of user+pass authentication, this tool uses OAUTH2 which is way more secure.
Unencrypted Google password may be leaked. In opposite, storing just an auth token is no big deal.

This tool was written mostly for demonstrational purposes, 
but actually I use this for my [i3blocks](https://github.com/vivien/i3blocks/) mail checker.

## Installation

    go get github.com/hypnoglow/gomuche

## Usage

Go to [Google Developers Console](https://console.developers.google.com/) and create your credentials:
client ID and client secret.

Create JSON file `$HOME/.gomuche/config.json` with the following content:

    {
        "client_id": "YOUR CLIENT ID",
        "client_secret": "YOUR CLIENT SECRET"
    }

Run:

    gomuche auth
    
It will print the link you need to open in your browser to authorize. You will be given a code.
Copy it and run mail check passing this code:

    gomuche check -c YOUR_CODE
    
It creates your auth token, which will be stored in `$HOME/.gomuche/token.json`. Next time you should
run mail check without passing the code:

    gomuche check
    
It will use your stored token, and automatically refresh it when it expires.

Note that you cannot use the same code again.
You must obtain a new one from `gomuche auth` to make a new token if you lose yours.

For any additional info, run:

    gomuche --help

## Issues

If you experience any problems (e.g. program exits with non-zero exit code), 
you can add `-v` flag to `check` command to print output verbosely. Also, you can check error log,
which is located in `$HOME/.gomuche/gomuche.log`.

If you found any bug, feel free to create a GitHub issue or/and submit a pull request.

## License

[MIT](https://github.com/hypnoglow/gomuche/blob/master/LICENCE.md)
