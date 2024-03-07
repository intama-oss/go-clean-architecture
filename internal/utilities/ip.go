package utilities

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type ipResponse struct {
	IP string `json:"ip"`
}

func GetPublicIP(ipCheckURL string) (string, error) {
	agent := fiber.AcquireAgent()
	req := agent.Request()
	req.Header.Set("User-Agent", "go-clean-architecture")
	req.SetRequestURI(ipCheckURL + "/?format=json")

	if err := agent.Parse(); err != nil {
		return "", errors.New("failed to parse request")
	}

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return "", errors.New("failed to dial destination server")
	}

	if statusCode != 200 {
		return "", errors.New("failed to get public IP")
	}

	var response ipResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", errors.New("failed to unmarshal response")
	}

	if response.IP == "" {
		return "", errors.New("failed to get public IP")
	}

	return response.IP, nil
}
