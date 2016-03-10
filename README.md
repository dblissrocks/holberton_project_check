# Holberton Project Check


## What is this?

Holberton Project Check helps Holberton School's students make sure that all their files are the ones that are expected before they get QA reviewed by HolbertonCloud.
- By default, the program creates the directories and files required for each task.
- Use the flag '-check' when running the program to check if the required files exist.

You need an account on Holberton School's intranet in order to use this.


## How to use on Mac OS

Download [the latest release](https://github.com/rudyrigot/holberton_project_check/releases).

Then, change your directory to your project's directory, and run:
```
~/Downloads/holberton_project_check [-check]
```
(or wherever the binary was downloaded)


## Run from source

Install Go (don't forget to setup a GOPATH), clone the repository under `$GOPATH/src`, and then:
```
go get github.com/howeyc/gopass
go run *.go
```

## Build:

```
go build -o holberton_project_check *.go
```


## Contribution

Contributions are welcome, whether to fix some Go style (I am fairly new to Go), or to build more features. Just open a pull request!
