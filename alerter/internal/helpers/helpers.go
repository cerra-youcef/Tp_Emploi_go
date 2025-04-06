package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Config struct {
	ConfigURL    string
	TimetableURL string
	MailToken    string
	ApiURL       string
}

func FetchAlertsByResource(apiURL string, resource int) ([]string, error) {
	// Convert resource ID to string for the URL
	resp, err := http.Get(fmt.Sprintf("%s/alerts?ucaID=%d", apiURL, resource))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a slice of structs containing just the email
	var alerts []struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &alerts); err != nil {
		return nil, err
	}

	// Extract just the emails into a string slice
	emails := make([]string, len(alerts))
	for i, alert := range alerts {
		emails[i] = alert.Email
	}

	return emails, nil
}
