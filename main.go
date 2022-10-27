package main

import (
	"log"
	"net/http"

	runtime "github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

var count int

func handleRequest() (Response, error) {
	log.Println("start handler")
	defer log.Println("end handler")

	count++
	message := "Hello AWS Lambda"
	for i := 0; i < count; i++ {
		message = message + "!"
	}

	log.Print(message)

	return Response{
		Message:    message,
		StatusCode: http.StatusOK,
	}, nil
}

func init() {
	count = 0
	log.Println("init function called")
}

func main() {
	runtime.Start(handleRequest)
}
