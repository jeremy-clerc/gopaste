# Gopaste

Some sort of pastebin like, very very simple, does not require any database, 
pastes are written to the filesystem, syntax highliting is done with 
[SyntaxHighlighter 3](http://alexgorbatchev.com/SyntaxHighlighter/).

## Features
* Only Syntax Highliting for now

## Run

Developped using go1.3 other versions may or may not work.

```
git clone https://github.com/jeremy-clerc/gopaste.git
cd gopaste
go run
```


By default it listen on port 8080

This version works, there is no input size limits, submit rate limit, so be
careful where you use it.

## Clean expired pastes

Use maintenance/remove_expired.pl in cron job

Default expiration is 1hour, max expiration is 1 week

## Todo
* Add more test
* Add screenshot
* Add input limit (2MB for example)
* Add password protection
* Create configuration file or cmd line arguments
* Create perl syntax highlighter for http://prismjs.com/ and migrate to prismjs

