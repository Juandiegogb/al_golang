package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	_ "github.com/lib/pq"
)

type secretValue struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
}

func getSecretValue(ctx context.Context) (*secretValue, error) {
	secretName := "prod/manager_db"
	region := "us-east-1"

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret: %w", err)
	}

	var secret secretValue
	if err := json.Unmarshal([]byte(*result.SecretString), &secret); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &secret, nil
}

func initDB() (*sql.DB, error) {
	ctx := context.Background()
	secret, err := getSecretValue(ctx)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		secret.Host, secret.Port, secret.Username, secret.Password, secret.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

type Activity struct {
	Id           int    `json:"id"`
	ActivityType string `json:"activity_type"`
	UnitId       int    `json:"unid_id"`
}

var db *sql.DB

func init() {

}

func handler(r events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	unitId, ok := r.QueryStringParameters["unit_id"]
	if !ok {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: "unit_id is required"}, nil
	}
	db, err := initDB()
	if err != nil {
		log.Fatalf("Application startup failed: %v", err)
	}
	defer db.Close()
	rows, err := db.Query("select sa.id,sa.activity_type,ap.unit_id from app_scheduledactivity sa left join app_activityplan ap on sa.activity_plan_id = ap.id where ap.unit_id = $1", unitId)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}
	var countries []Activity
	for rows.Next() {
		var c Activity
		if err := rows.Scan(&c.Id, &c.ActivityType, &c.UnitId); err != nil {
			return events.APIGatewayV2HTTPResponse{}, err
		}
		countries = append(countries, c)
	}
	if err := rows.Close(); err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: "Internal error"}, err
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(countries)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: "Internal error"}, err
	}
	return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: string(data)}, nil
}

func main() {
	lambda.Start(handler)
}
