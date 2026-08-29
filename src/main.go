package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var response events.APIGatewayProxyResponse

	switch r.Path {
	case "/api/greet":
		response = events.APIGatewayProxyResponse{StatusCode: 200, Body: "hello"}
	case "api/data":
		person := Person{Name: "Juan Diego", Age: 22}
		encoded_json, err := json.Marshal(person)
		if err != nil {
			response = events.APIGatewayProxyResponse{StatusCode: 500}
			break
		}
		response = events.APIGatewayProxyResponse{StatusCode: 200, Body: string(encoded_json)}
	default:
		response = events.APIGatewayProxyResponse{StatusCode: 404, Body: "Not found"}
	}
	fmt.Println(r)
	return response, nil
}

func main() {
	lambda.Start(handler)
}
