package todos

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	uuid "github.com/satori/go.uuid"
)

const TableName = "Todos"

func createUUID() string {
	uuid := uuid.NewV4()
	return uuid.String()
}

type Todo struct {
	ID          string `json:"id"`
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
}

func GetTodos() ([]Todo, error) {
	var todos []Todo
	svc := dynamodb.New(session.New())

	input := &dynamodb.ScanInput{
		TableName: aws.String(TableName),
	}

	result, err := svc.Scan(input)
	if err != nil {
		fmt.Errorf("failed to make Query API call, %v", err)
		return todos, err
	}

	dynamodbattribute.UnmarshalListOfMaps(result.Items, &todos)
	if err != nil {
		fmt.Errorf("failed to unmarshal Query result items, %v", err)
		return todos, err
	}

	return todos, nil
}

func GetTodo(id string) (Todo, error) {
	todo := Todo{}
	svc := dynamodb.New(session.New())

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(TableName),
	}

	result, err := svc.GetItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return todo, err
	}

	// Unmarshall the result in to a Todo
	err = dynamodbattribute.UnmarshalMap(result.Item, &todo)

	if err != nil {
		fmt.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

func DeleteTodo(id string) error {
	svc := dynamodb.New(session.New())

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func PostTodo(description string) (Todo, error) {
	svc := dynamodb.New(session.New())
	id := createUUID()

	fmt.Println("id", id)

	// Marshall the request body
	todo := Todo{
		Description: description,
		ID:          id,
		Completed:   false,
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Description": {
				S: aws.String(todo.Description),
			},
			"Id": {
				S: aws.String(todo.ID),
			},
			"Completed": {
				BOOL: aws.Bool(todo.Completed),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(TableName),
	}

	_, err := svc.PutItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return todo, err
	}

	return todo, nil
}

func CompleteTodo(id string, completed bool) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":b": {
				BOOL: aws.Bool(completed),
			},
		},
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set info.Completed = :b"),
	}

	_, err := svc.UpdateItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
