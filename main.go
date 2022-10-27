package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	runtime "github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	RequestID  string `json:"requestId"`
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
	log.Printf("Message: %s", message)

	client := new(http.Client)
	requestID, err := getRequestID(client)
	if err != nil {
		log.Print(err.Error())
	}
	log.Printf("RequestID: %s", requestID)

	return Response{
		RequestID:  requestID,
		Message:    message,
		StatusCode: http.StatusOK,
	}, nil
}

const (
	runtimeAPIEnvKey           = "AWS_LAMBDA_RUNTIME_API"
	runtimeRequestIDHeaderName = "Lambda-Runtime-Aws-Request-Id"
)

func getRequestID(client *http.Client) (string, error) {
	host := os.Getenv(runtimeAPIEnvKey)
	if host == "" {
		return "", fmt.Errorf("host is empty")
	}

	url := fmt.Sprintf("http://%s/2018-06-01/runtime/invocation/next", host)
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	rHeader := res.Header
	rIds, exists := rHeader[runtimeRequestIDHeaderName]
	if !exists {
		return "", fmt.Errorf("'%s' header does not exists.", runtimeRequestIDHeaderName)
	}

	if len(rIds) == 0 {
		return "", fmt.Errorf("Value of '%s' header is empty.", runtimeRequestIDHeaderName)
	}

	return rIds[0], nil
}

func init() {
	count = 0
	log.Println("init function called")
}

func main() {
	runtime.Start(handleRequest)
}
