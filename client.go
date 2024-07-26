package opengo

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type Client struct {
	config Config
}

type Config struct {
	authToken string
}

func CreateClient(authToken string) *Client {
	config := Config{
		authToken: authToken,
	}
	return &Client{
		config: config,
	}
}

func (c *Client) sendRequest(method string, requestBody *bytes.Buffer, writer *multipart.Writer, url string) (string, error) {
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.config.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(bodyBytes), nil
}
