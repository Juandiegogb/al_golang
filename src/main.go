package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(r events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var response events.APIGatewayV2HTTPResponse

	switch r.RawPath {
	case "/api/greet":
		response = events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: "hello"}
	case "/api/data":
		person := Person{Name: "Juan Diego", Age: 22}
		encoded_json, err := json.Marshal(person)
		if err != nil {
			response = events.APIGatewayV2HTTPResponse{StatusCode: 500}
			break
		}
		response = events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: string(encoded_json)}
	default:
		response = events.APIGatewayV2HTTPResponse{StatusCode: 404, Body: "Not found"}
	}
	log.Println(r.RawPath)
	return response, nil
}

func main() {
	lambda.Start(handler)
}
