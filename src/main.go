package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Person struct {
	Name           string `json:"name"`
	Age            int    `json:"age"`
	AdditionalInfo string `json:"additional_info"`
}

func handleDynamoUpdate() (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: "Hi, function is fine"}, nil
}

func handleGetData(unitId string) (events.APIGatewayV2HTTPResponse, error) {
	person := Person{Name: "Juan Diego", Age: 22, AdditionalInfo: unitId}
	encoded_json, err := json.Marshal(person)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, nil
	}
	return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: string(encoded_json)}, nil
}

func handler(r events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	var (
		response events.APIGatewayV2HTTPResponse
		err      error
	)

	switch r.RawPath {
	case "/api/greet":
		response, err = handleDynamoUpdate()
	case "/api/data":
		unitId, ok := r.QueryStringParameters["unit_id"]
		if !ok {
			response = events.APIGatewayV2HTTPResponse{StatusCode: http.StatusBadRequest, Body: "unit_id not found"}
			break
		}
		response, err = handleGetData(unitId)
	default:
		response = events.APIGatewayV2HTTPResponse{StatusCode: 404, Body: "Not found"}
	}
	log.Println(r.RawPath)
	return response, err
}

func main() {
	lambda.Start(handler)
}
