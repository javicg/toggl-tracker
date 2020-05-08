package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var ENV_USERNAME = "TOGGL_USERNAME"
var ENV_PASSWORD = "TOGGL_PASSWORD"

type TogglApi struct {
	baseUrl string
	client  *http.Client
}

func NewTogglApi() (api *TogglApi) {
	api = &TogglApi{}
	api.baseUrl = "https://www.toggl.com/api/v8"
	api.client = &http.Client{}
	return api
}

type Me struct {
	Data struct {
		Email    string
		Fullname string
	}
}

func (toggl *TogglApi) GetMe() (me Me, err error) {
	resp, err := toggl.getAuthenticated("/me")
	if err != nil {
		fmt.Println("[GetMe] Request failed! Error:", err)
		return
	} else if resp.StatusCode != 200 {
		fmt.Println("[GetMe] Request failed with status:", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&me)
	if err != nil {
		fmt.Println("[GetMe] Error unmarshalling response:", err)
		return
	}

	return
}

type TimeEntry struct {
	Id          int
	Pid         int
	Start       time.Time
	Stop        time.Time
	Duration    int
	Description string
	Tags        []string
}

func (toggl *TogglApi) GetTimeEntries(start time.Time, end time.Time) (entries []TimeEntry, err error) {
	params := map[string]string{
		"start_date": start.Format(time.RFC3339),
		"end_date":   end.Format(time.RFC3339),
	}

	resp, err := toggl.getAuthenticatedWithQueryParams("/time_entries", params)
	if err != nil {
		fmt.Println("[GetTimeEntries] Request failed! Error:", err)
		return
	} else if resp.StatusCode != 200 {
		fmt.Println("[GetTimeEntries] Request failed with status:", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		fmt.Println("[GetTimeEntries] Error unmarshalling response:", err)
		return
	}
	return
}

func (toggl *TogglApi) getAuthenticatedWithQueryParams(path string, params map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", toggl.baseUrl+path, nil)
	if err != nil {
		return
	}

	err = addBasicAuth(req)
	if err != nil {
		return
	}

	q := req.URL.Query()
	for p := range params {
		q.Add(p, params[p])
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")
	return toggl.client.Do(req)
}

func (toggl *TogglApi) getAuthenticated(path string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", toggl.baseUrl+path, nil)
	if err != nil {
		return
	}

	err = addBasicAuth(req)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	return toggl.client.Do(req)
}

func addBasicAuth(req *http.Request) (err error) {
	username, ok := os.LookupEnv(ENV_USERNAME)
	if !ok {
		err = missingEnv(ENV_USERNAME)
		return
	}

	password, ok := os.LookupEnv(ENV_PASSWORD)
	if !ok {
		err = missingEnv(ENV_PASSWORD)
		return
	}

	req.SetBasicAuth(username, password)
	return
}

func missingEnv(env string) error {
	return errors.New(fmt.Sprintf("%s not specified!", env))
}
