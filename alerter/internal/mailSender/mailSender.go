package mailSender

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type EmailRequest struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
}

func SendEmail(to, subject, content, token, apiURL string) error {

	email := EmailRequest{
		Recipient: to,
		Subject:   subject,
		Content:   content,
	}

	jsonData, err := json.Marshal(email)
	if err != nil {
		return errors.New("error in json serialization: " + err.Error())
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("error while creating request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token) // pas de "Bearer" ici selon la doc

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("Error while sending query: " + err.Error())

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New("API Error: " + resp.Status + string(body))
	}

	return nil
}
