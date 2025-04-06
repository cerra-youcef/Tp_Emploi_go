package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	ConfigURL    string
	TimetableURL string
	MailToken    string
	ApiURL       string
}

// Fetch alerts from Config to get emails
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

func LoadConfig() (Config, error) {
	var cfg Config
	var ok bool

	if cfg.ConfigURL, ok = os.LookupEnv("CONFIG_URL"); !ok {
		return cfg, errors.New("CONFIG_URL not set in .env file")
	}

	if cfg.TimetableURL, ok = os.LookupEnv("TIMETABLE_URL"); !ok {
		return cfg, errors.New("TIMETABLE_URL not set in .env file")
	}

	if cfg.MailToken, ok = os.LookupEnv("MAIL_TOKEN"); !ok {
		return cfg, errors.New("MAIL_TOKEN not set in .env file")
	}

	if cfg.ApiURL, ok = os.LookupEnv("API_URL"); !ok {
		return cfg, errors.New("MAIL_TOKEN not set in .env file")
	}
	return cfg, nil
}
