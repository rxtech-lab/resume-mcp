package utils

import (
	"fmt"
	"net/url"
	"os"
)

func GetTransactionSessionUrl(serverPort string, sessionId string) (string, error) {
	// Override baseUrl if BASE_URL env var is set
	if os.Getenv("BASE_URL") != "" {
		baseUrl := os.Getenv("BASE_URL")
		parsedUrl, err := url.Parse(baseUrl)
		if err != nil {
			return "", fmt.Errorf("invalid BASE_URL env var: %w", err)
		}
		parsedUrl.Path = fmt.Sprintf("/resume/preview/%s", sessionId)
		return parsedUrl.String(), nil
	}

	url := fmt.Sprintf("http://localhost:%s/resume/preview/%s", serverPort, sessionId)
	return url, nil
}

func GetDownloadSessionUrl(serverPort string, sessionId string) (string, error) {
	// Override baseUrl if BASE_URL env var is set
	if os.Getenv("BASE_URL") != "" {
		baseUrl := os.Getenv("BASE_URL")
		parsedUrl, err := url.Parse(baseUrl)
		if err != nil {
			return "", fmt.Errorf("invalid BASE_URL env var: %w", err)
		}
		parsedUrl.Path = fmt.Sprintf("/resume/download/%s", sessionId)
		return parsedUrl.String(), nil
	}

	url := fmt.Sprintf("http://localhost:%s/resume/download/%s", serverPort, sessionId)
	return url, nil
}
