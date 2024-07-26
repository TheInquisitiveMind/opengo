package main

import (
	"errors"
	"flag"
	"fmt"
	"opengo"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

const validExtensions = ".mp3|.wav|.flac|.aac|.ogg|.m4a"

func main() {
	/* Usage:
	go run speechtotext.go -file=sound.mp3
	go run speechtotext.go -file=sound.mp3 --temperature=1.0
	*/

	if len(os.Args) < 2 {
		fmt.Println("Please provide a path to file\nUsage: go run ./speechtotext <file_path>")
		return
	}

	filePath := flag.String("file", "", "Path to the audio file (required)")
	temperature := flag.Float64("temperature", 0.0, "Model temperature value (optional) <0.0;2.0>")
	flag.Parse()

	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}
	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set")
		return
	}

	if err := validateFile(*filePath); err != nil {
		fmt.Printf("Invalid file: %v\n", err)
		return
	}

	if err := validateTemperature(*temperature); err != nil {
		fmt.Printf("Invalid temperature: %v\n", err)
		os.Exit(1)
	}

	client := opengo.CreateClient(apiKey)
	resp, err := client.ChangeToText(
		opengo.SoundRequestParams{
			FilePath:    *filePath,
			Temperature: *temperature,
		},
	)

	if err != nil {
		fmt.Printf("Error during transcription: %v\n", err)
		return
	}
	fmt.Println(resp)
}

func validateTemperature(temp float64) error {
	if temp < 0.0 || temp > 2.0 {
		return errors.New("temperature must be between 0.0 and 2.0")
	}
	return nil
}

func validateFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("could not access file: %w", err)
	}
	fileExt := strings.ToLower(filepath.Ext(filePath))
	if !strings.Contains(validExtensions, fileExt) {
		return fmt.Errorf("file is not a supported audio format")
	}
	return nil
}
