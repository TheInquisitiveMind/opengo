package opengo

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

const modelName = "whisper-1"
const url = "https://api.openai.com/v1/audio/transcriptions"

type SoundRequestParams struct {
	FilePath    string
	Temperature float64
}

func (c *Client) ChangeToText(request SoundRequestParams) (string, error) {
	file, err := os.Open(request.FilePath)
	if err != nil {
		return "", fmt.Errorf("error opening audio file: %v", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", request.FilePath)

	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copying file to form file: %v", err)
	}

	err = writer.WriteField("model", modelName)
	if err != nil {
		return "", fmt.Errorf("set of model: %w", err)
	}
	if request.Temperature != 0 {
		err = writer.WriteField("temperature", fmt.Sprintf("%.2f", request.Temperature))
		if err != nil {
			return "", fmt.Errorf("set of temperature: %w", err)
		}
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing writer: %v", err)
	}

	return c.sendRequest("POST", &requestBody, writer, url)
}
