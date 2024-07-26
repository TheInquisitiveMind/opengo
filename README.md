# opengo
OpenAI client written in Go

before run:
1. go mod download github.com/joho/godotenv

2. cp .env_example .env

3. set your api key in .env



Example of usage:

go run speechtotext.go -file=sound.mp3
go run speechtotext.go -file=sound.mp3 --temperature=1.0

