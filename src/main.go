package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(r)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "hola"}, nil
}

func main() {
	lambda.Start(handler)
}
