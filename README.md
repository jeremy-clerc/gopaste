# Gopaste

[![Build Status](https://travis-ci.org/jeremy-clerc/gopaste.svg?branch=master)](https://travis-ci.org/jeremy-clerc/gopaste)

Some sort of pastebin like, very very simple, does not require any database, 
pastes are written to the filesystem, syntax highliting is done with 
[SyntaxHighlighter 3](http://alexgorbatchev.com/SyntaxHighlighter/).

![Paste it!](http://i.imgur.com/5jlZnlV.png "Paste it!")
![Share it!](http://i.imgur.com/CgUtLcg.png "Share it!")
![View it!](http://i.imgur.com/BbbFGk4.png "View it!")

## Features
* Syntax Highliting
* Choose expiration time between 1hour to 1week
* Paste size limit is 1MB
* Url is dynamic, it uses the one you are accessing

## Run

Developped using go1.3 other versions may or may not work.

```
git clone https://github.com/jeremy-clerc/gopaste.git
cd gopaste
go run
```


By default it listen on port 8080

## Clean expired pastes

Use maintenance/remove_expired.pl in cron job

Default expiration is 1hour, max expiration is 1 week

## Todo
* Add more test
* Add password protection
* Create configuration file or cmd line arguments
* Create perl syntax highlighter for http://prismjs.com/ and migrate to prismjs

