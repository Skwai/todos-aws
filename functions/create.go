package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/skwai/todos-aws/todos"
)

type TodoData struct {
	Description string `json:"description"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data := TodoData{}
	err := json.Unmarshal([]byte(request.Body), &data)

	if err != nil {
		fmt.Println(err.Error)
	}

	fmt.Println("data", data)
	fmt.Println("data.Description", data.Description)
	fmt.Println("body", request.Body)

	todo, err := todos.PostTodo(data.Description)

	if err != nil {
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: "Error", StatusCode: 500}, nil
	}

	j, _ := json.Marshal(todo)
	s := string(j)
	return events.APIGatewayProxyResponse{Body: s, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
