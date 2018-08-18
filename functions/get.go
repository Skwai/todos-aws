package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/skwai/todos-aws/todos"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	todo, err := todos.GetTodo(id)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	if todo.ID == "" {
		return events.APIGatewayProxyResponse{Body: "Not Found", StatusCode: 404}, nil
	}

	j, _ := json.Marshal(todo)
	s := string(j)
	return events.APIGatewayProxyResponse{Body: s, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
