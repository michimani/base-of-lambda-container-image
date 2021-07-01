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

func handleRequest() (Response, error) {
	log.Println("start handler")
	defer log.Println("end handler")

	return Response{
		Message:    "Hello AWS Lambda",
		StatusCode: http.StatusOK,
	}, nil
}

func init() {
	log.Println("init function called")
}

func main() {
	runtime.Start(handleRequest)
}
