package utilities

import (
	"encoding/json"
	"errors"
	"net/http"
)

func GetPublicIP() (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.ipify.org?format=json", nil)

	req.Header.Set("User-Agent", "go-clean-architecture")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}

	if ip, ok := body["ip"].(string); ok {
		return ip, nil
	}

	return "", errors.New("failed to get public IP")
}
