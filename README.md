# Holberton Project Check


## What is this?

Holberton Project Check helps Holberton School's students make sure that all their files are the ones that are expected before they get QA reviewed by HolbertonCloud.

It is composed of two main files:
 * `holberton_intranet_auth.go` provides all the token-based authentication to the intranet.
 * `main.go` will ultimately provide the end-user features; for now is pretty minimal.

You need an account on Holberton School's intranet in order to use this.


## About the token-based authentication

The first time the user logs in, a configuration file is created in the user's repository containing a temporary token.
From there on, all of the user cycle is managed; for instance, if the server refuses the token, the user is cleanly prompted to log in again.


## Install & run locally

 * Install:
```
go get github.com/howeyc/gopass
```
 * Run:
```
go run *.go
```
 * Build:
```
go build *.go
```


## Contribution

Contributions are welcome, whether to fix some Go style (I am fairly new to Go), or to build more features. Just open a pull request!
