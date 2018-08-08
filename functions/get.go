package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/skwai/todos-aws/todos"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	todo, err := todos.GetTodo(id)

	if err != nil {
		panic(fmt.Sprintf("Failed to find Item, %v", err))
	}

	// Log and return result
	j, _ := json.Marshal(todo)
	s := string(j)
	fmt.Println("Found item: ", s)
	return events.APIGatewayProxyResponse{Body: s, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
