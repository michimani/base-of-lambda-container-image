package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Response struct {
	RequestIDFromRuntimeAPI string `json:"requestIdFromRuntimeAPI"`
	RequestIDFromContext    string `json:"requestIdFromContext"`
	Message                 string `json:"message"`
	StatusCode              int    `json:"statusCode"`
}

var count int

func handleRequest(ctx context.Context) (Response, error) {
	log.Println("start handler")
	defer log.Println("end handler")

	count++
	message := "Hello AWS Lambda"
	for i := 0; i < count; i++ {
		message = message + "!"
	}
	log.Printf("Message: %s", message)

	// Get AWS Request ID from Lambda Runtime API.
	client := new(http.Client)
	requestIDFromRuntimeAPI, err := getRequestID(client)
	if err != nil {
		log.Print(err.Error())
	}
	log.Printf("RequestID from Runtime API: %s", requestIDFromRuntimeAPI)

	// Get AWS Request ID from context.
	requestIDFromContext := ""
	lc, exists := lambdacontext.FromContext(ctx)
	if !exists {
		log.Print("Failed to get Lambda Context from context.")
	} else {
		requestIDFromContext = lc.AwsRequestID
	}
	log.Printf("RequestID from context: %s", requestIDFromContext)

	return Response{
		RequestIDFromRuntimeAPI: requestIDFromRuntimeAPI,
		RequestIDFromContext:    requestIDFromContext,
		Message:                 message,
		StatusCode:              http.StatusOK,
	}, nil
}

const (
	runtimeAPIEnvKey           = "AWS_LAMBDA_RUNTIME_API"
	runtimeRequestIDHeaderName = "Lambda-Runtime-Aws-Request-Id"
)

// getRequestID returns AWS Request ID from Lambda Runtime API
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
