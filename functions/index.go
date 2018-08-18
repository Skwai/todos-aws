package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/skwai/todos-aws/todos"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	todos, err := todos.GetTodos()

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	j, _ := json.Marshal(todos)
	s := string(j)
	return events.APIGatewayProxyResponse{Body: s, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
