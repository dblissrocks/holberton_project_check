package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"

	"github.com/howeyc/gopass"
)

const holbertonDirName string = ".holberton"
const holbertonIntranetTokenFile string = "intranet_token"

// To unmarshal the responses from the server about tokens
type tokenAuthFromServer struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

// To marshal and unmarshal to and from the config file
type tokenAuth struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// Will get the URL, handle all of the Holberton intranet auth, and return the body
func getWithHolbertonAuth(urlToGet string) (string, error) {
	auth, err := getEmailAndToken()
	if err != nil {
		return "", err
	}

	urlToGetWithParams := fmt.Sprintf("%s?user_email=%s&user_token=%s", urlToGet, auth.Email, auth.Token)

	resp, err := http.Get(urlToGetWithParams)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 401 { // Unauthorized = try again
		fmt.Println("Authentication was refused, probably because your token expired.")
		// login (will overwrite the config file), and make the call again
		if _, err := login(); err != nil {
			return "", err
		}
		return getWithHolbertonAuth(urlToGet)
	}

	// 200 = let's proceed
	if resp.StatusCode == 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return buf.String(), nil // we did not close resp.Body
	}

	var m tokenAuthFromServer
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", err
	}
	// Other HTTP error code = get the error message and stop there
	return "", errors.New(m.Message)
}

// Will get it one way or another: will try getting the file first, fetch it later
// Returns the email of the user and its token
func getEmailAndToken() (tokenAuth, error) {
	holbertonDirPath, err := ensureHolbertonPersonalDirectory()
	if err != nil {
		return tokenAuth{}, err
	}
	holbertonIntranetTokenFilePath := path.Join(holbertonDirPath, holbertonIntranetTokenFile)
	if _, err := os.Stat(holbertonIntranetTokenFilePath); os.IsNotExist(err) { // File doesn't exist
		fmt.Println("You never authenticated on this computer before.")
		auth, err := login()
		if err != nil {
			return tokenAuth{}, err
		}
		return auth, nil // Returning the token we got from the login
	}

	// File exists -> read in the file
	fmt.Println("Auth token found on this computer.")

	configFile, err := os.Open(holbertonIntranetTokenFilePath)
	if err != nil {
		return tokenAuth{}, err
	}
	defer configFile.Close()

	var t tokenAuth

	if err := json.NewDecoder(configFile).Decode(&t); err != nil {
		return tokenAuth{}, err
	}

	return t, nil
}

// Ensures the config repository (~/.holberton) exists, creates it otherwise
func ensureHolbertonPersonalDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	holbertonDirPath := path.Join(usr.HomeDir, holbertonDirName)
	if _, err := os.Stat(holbertonDirPath); os.IsNotExist(err) {
		if err := os.Mkdir(holbertonDirPath, 0700); err != nil {
			return "", err
		}
	}
	return holbertonDirPath, nil
}

// Will login, fetch the token from the intranet, write it to the config file with the user's email, and return them both
func login() (tokenAuth, error) {

	// Prompting the user for login info
	var email string
	var password []byte

	fmt.Println("Please login with your intranet's credentials.")
	fmt.Print("Email: ")
	fmt.Scanln(&email)

	fmt.Print("Password: ")
	password, err := gopass.GetPasswd()
	if err != nil {
		return tokenAuth{}, err
	}

	// Fetch the token from the endpoint
	resp, err := http.PostForm("https://intranet.hbtn.io/token", url.Values{"email": {email}, "password": {string(password)}})
	if err != nil {
		return tokenAuth{}, err
	}
	defer resp.Body.Close()

	var m tokenAuthFromServer
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return tokenAuth{}, err
	}

	// Deal with errors informed by the server (Invalid credentials, ...)
	if !m.Success {
		return tokenAuth{}, errors.New(m.Message)
	}
	token := m.Token

	// We want to store the email AND the token
	tokenToWrite, err := json.Marshal(tokenAuth{Email: email, Token: token})
	if err != nil {
		return tokenAuth{}, err
	}

	// Let's make sure the config directory exists first, and get the path we want
	holbertonDirPath, err := ensureHolbertonPersonalDirectory()
	if err != nil {
		return tokenAuth{}, err
	}
	holbertonIntranetTokenFilePath := path.Join(holbertonDirPath, holbertonIntranetTokenFile)
	if err := ioutil.WriteFile(holbertonIntranetTokenFilePath, []byte(tokenToWrite), 0700); err != nil {
		return tokenAuth{}, err
	}

	return tokenAuth{Email: email, Token: token}, nil
}
