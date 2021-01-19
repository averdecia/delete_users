package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	command "github.com/averdecia/script_command"
)

// RemoveUsers struct is created for implementing ICommand interface
type RemoveUsers struct {
	args Args
}

// ExecuteAction implements ICommand interface
func (c *RemoveUsers) ExecuteAction(element []string) (string, error) {
	fmt.Printf("Element: %v", element)
	email := element[1]
	gamification, err := c.getUserGamification(email)
	if err != nil {
		return "Fail", err
	}
	status, err := c.removeUser(gamification)

	if err != nil {
		return "Fail", err
	}

	return status, nil
}

func (c *RemoveUsers) removeUser(gamification string) (string, error) {
	client := http.Client{
		Timeout: time.Duration(5 * time.Minute),
	}

	params := url.Values{
		"gamificationid": {gamification},
	}

	request, err := http.NewRequest("POST", c.args.Endpoint+"/users/remove", strings.NewReader(params.Encode()))
	request.Header.Set("Authorization", "Basic "+c.args.AuthToken)
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("Cannot connect to server: %v", err)
		return "Error", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.Bytes()

	r := bytes.NewReader(body)
	decoder := json.NewDecoder(r)

	response := &LookupResponse{}
	err = decoder.Decode(response)

	if err != nil {
		return "Error", err
	}

	return response.Status, nil
}

func (c *RemoveUsers) getUserGamification(email string) (string, error) {
	client := http.Client{
		Timeout: time.Duration(5 * time.Minute),
	}

	request, err := http.NewRequest("GET", c.args.Middleware+
		"/services/user/checkemail?email="+email+"&country_code=MX", nil)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("Cannot connect to server: %v", err)
		return "Error", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.Bytes()

	r := bytes.NewReader(body)
	decoder := json.NewDecoder(r)

	response := &MWresponse{}
	err = decoder.Decode(response)

	if err != nil {
		return "Error", err
	}

	if !response.Data.Response.Exist {
		return "Error", errors.New("User not found")
	}
	gamification := ""
	for _, s := range response.Data.Response.SocialData {
		if s.RedSocial == "IMUSICA" {
			gamification = s.GamificationID
		}
	}
	if gamification == "" {
		return "Error", errors.New("User not found")
	}

	return gamification, nil
}

func getArgs() Args {
	argsOS := os.Args[1:]
	routines, _ := strconv.Atoi(argsOS[3])
	return Args{
		Endpoint:   argsOS[0],
		InputPath:  argsOS[1],
		OutputPath: argsOS[2],
		GoRoutines: routines,
		AuthToken:  argsOS[4],
		Middleware: argsOS[5],
	}
}

func main() {
	// You can use os.Args[1:] to pass the variables to the build
	args := getArgs()
	mycommand := &RemoveUsers{
		args: args,
	}
	command.RunProcess(mycommand, mycommand.args.GoRoutines, mycommand.args.InputPath, mycommand.args.OutputPath)
}
